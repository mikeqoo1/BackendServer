package middleware

import (
	"fmt"

	"errors"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var (
	//ErrTokenExpired Token過期
	ErrTokenExpired error = errors.New("Token is expired")
	//ErrTokenNotValidYet Token not active yet
	ErrTokenNotValidYet error = errors.New("Token not active yet")
	//ErrTokenMalformed That's not even a token
	ErrTokenMalformed error = errors.New("That's not even a token")
	//ErrTokenInvalid Couldn't handle this token
	ErrTokenInvalid error = errors.New("Couldn't handle this token")
	//SignKey 簽名訊息, 應該要從資料庫來取得
	SignKey string = "Seaflower"
)

//JWTAuth 中間件, 檢查token
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		if token == "" {
			c.JSON(http.StatusOK, gin.H{
				"status": -1,
				"msg":    "請求未攜帶token， 無許可權訪問",
				"data":   nil,
			})
			c.Abort()
			return
		}

		log.Print("get token: ", token)
		j := NewJWT()
		//解析token中包含的相關資料
		claims, err := j.ParserToken(token)

		fmt.Println(claims, err)
		if err != nil {
			//token過期
			if err == ErrTokenExpired {
				c.JSON(http.StatusOK, gin.H{
					"status": -1,
					"msg":    "token授權已過期， 請重新申請授權",
					"data":   nil,
				})
				c.Abort()
				return
			}
			//其他Error
			c.JSON(http.StatusOK, gin.H{
				"status": -1,
				"msg":    err.Error(),
				"data":   nil,
			})
			c.Abort()
			return
		}

		//解析到具體的claims相關資料
		c.Set("claims", claims)

	}
}

//JWT JWT基本結構, 簽名的signkey
type JWT struct {
	SigningKey []byte
}

//CustomClaims 定義要求的內容
type CustomClaims struct {
	//UserAccount Account帳號
	UserAccount string `json:"account"`
	//StandardClaims結構實現了Claims接口(Valid()函数)
	jwt.StandardClaims
}

//NewJWT init JWT
func NewJWT() *JWT {
	return &JWT{
		[]byte(GetSignKey()),
	}
}

//GetSignKey 取得signkey
func GetSignKey() string {
	return SignKey
}

//SetSignKey 設定signkey
func SetSignKey(key string) string {
	SignKey = key
	return SignKey
}

//CreateToken 使用User基本資料claims以及簽名key(signkey)生成token
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	//https://gowalker.org/github.com/dgrijalva/jwt-go#Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

//ParserToken token解析
func (j *JWT) ParserToken(tokenString string) (*CustomClaims, error) {
	//https://gowalker.org/github.com/dgrijalva/jwt-go#ParseWithClaims
	//輸入User自定的Claims struct, token, 以及自定義func來解析token字串為jwt的Token指針
	//Keyfunc是匿名函數: type Keyfunc func(*Token) (interface{}, error)
	//func ParseWithClaims(tokenString string, claims Claims, keyFunc Keyfunc) (*Token, error) {}
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})

	fmt.Println(token, err)
	if err != nil {
		//https://gowalker.org/github.com/dgrijalva/jwt-go#ValidationError
		//jwt.ValidationError 是無效token的error struct
		if ve, ok := err.(*jwt.ValidationError); ok {
			//ValidationErrorMalformed是一个uint常量，表示token不可用
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, ErrTokenMalformed
				//ValidationErrorExpired表示Token過期
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, ErrTokenExpired
				//ValidationErrorNotValidYet表示無效的token
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, ErrTokenNotValidYet
			} else {
				return nil, ErrTokenInvalid
			}

		}
	}

	//把token中的claims資料解析出来和User原始資料進行比對
	//做類型判斷， 把token.Claims轉換成我們自定義的Claims結構
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrTokenInvalid

}

//UpdateToken 更新Token
func (j *JWT) UpdateToken(tokenString string) (string, error) {
	//TimeFunc是一個默認值是time.Now, 用來解析token後花進行過期時間驗證
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}

	//拿到token基本資料
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil

	})

	//比對token當下還有效
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		//修改Claims的過期時間(int64)
		//https://gowalker.org/github.com/dgrijalva/jwt-go#StandardClaims
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return j.CreateToken(*claims)
	}
	return "", fmt.Errorf("token獲取失敗:%v", err)
}
