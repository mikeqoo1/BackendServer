package router

import (
	mylib "BackendServer/mylib"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

//用來判斷是哪種資料 0代表raw data 1代表form-data
var checkflag int

//Userinfo 管理員
type Userinfo struct {
	Account  string `json:"account"`
	Password string `json:"Password"`
}

//Logininfo 登入結果
type Logininfo struct {
	Token    string `json:"token"`
	Account  string `json:"account"`
	Password string `json:"Password"`
}

//Registeruser 註冊
func Registeruser(c *gin.Context) {
	var info Userinfo
	err := c.BindJSON(&info)
	if err == nil {
		stmt, err := mylib.ConstDBpool.Prepare("INSERT INTO `資料表` SET `欄位1` = ?, `欄位2` = ?")
		if err != nil {
			mylib.MyLogger.Error("New SQL 語法失敗", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "後台sql有誤",
			})
			return
		}
		_, err1 := stmt.Exec(info.Account, info.Password)
		if err1 == nil {
			c.JSON(http.StatusOK, gin.H{
				"message": "成功註冊",
				"data":    nil,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"message": "註冊失敗" + err.Error(),
				"data":    nil,
			})
		}
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "user解析資料失敗" + err.Error(),
			"data":    nil,
		})
	}
}

//Login Login
func Login(c *gin.Context) {
	var login Logininfo
	if c.BindJSON(&login) == nil {
		//登錄的邏輯驗證
		isPass, user, err := loginCheck(login)
		if isPass {
			generateToken(c, user)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  -1,
				"message": "登入失敗" + err.Error(),
				"data":    nil,
			})
		}

	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "User資料解析失敗",
			"data":    nil,
		})
	}
}

//CheckData 檢查資料 中間件, 未來可以當驗證的檢查
func CheckData() gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := c.GetRawData()
		if err != nil {
			mylib.MyLogger.Error("CheckData Error: ", err.Error())
		}
		mylib.MyLogger.Info("原始 data: ", string(data))
		//fmt.Printf("原始 data: %v\n", string(data))

		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data))
		//fmt.Println(c.Request.Body)

		c.Next()
	}
}

//List List的api func
func List(c *gin.Context) {
	c.JSON(http.StatusOK, "hello word")
}

//New New的api func
func New(c *gin.Context) {
	checkflag = 0
	mylib.MyLogger.Debug("誰來新增: ", c.ClientIP())

	stmt, err := mylib.ConstDBpool.Prepare("INSERT INTO `資料表` SET `欄位1` = ?, `欄位2` = ?")
	if err != nil {
		mylib.MyLogger.Error("New SQL 語法失敗", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "後台sql有誤",
		})
		return
	}

	var info Userinfo

	//先用raw data解析
	err = c.BindJSON(&info)
	if err != nil {
		checkflag = 1
		mylib.MyLogger.Warn("該筆資訊是form-data")
	}

	if checkflag == 1 {
		user := c.DefaultPostForm("account", "NULL")
		if user == "NULL" {
			mylib.MyLogger.Error("請檢查格式或內容")
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "請檢查格式或內容",
			})
			return
		}
		info.Account = user
	}

	mylib.MyLogger.Debug("UserInfo=> ", info)

	_, err1 := stmt.Exec(info.Account, info.Password)
	//fmt.Println(i)
	if err1 != nil {
		mylib.MyLogger.Error("Create new user Error:", err1.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err1.Error(),
		})
	} else {
		mylib.MyLogger.Info("Create new user success")
		c.JSON(http.StatusOK, gin.H{
			"message": "Create new user success",
		})
	}
}

//Update Update的api func
func Update(c *gin.Context) {
	//參照上面範例
}

//Drop Drop的api func
func Drop(c *gin.Context) {
	//參照上面範例
}

func loginCheck(login Logininfo) (bool, Userinfo, error) {
	userData := Userinfo{}
	userExist := false
	var info Userinfo

	rows, err := mylib.ConstDBpool.Query("SELECT * FROM Userlist WHERE account = ?", info.Account)
	mylib.MyLogger.Debug("rows: ", rows)
	if err != nil {
		mylib.MyLogger.Error("Login SQL 語法失敗", err.Error())
		return userExist, userData, err
	}

	tablestruct, err := mylib.FromRows(rows)
	for _, ele := range tablestruct.Rows {
		info.Account = ele["account"].(string)
		info.Password = ele["password"].(string)
	}
	mylib.MyLogger.Debug("tablestruct info: ", info)

	if login.Account == info.Account && login.Password == info.Password {
		userExist = true
		userData.Account = info.Account
	}

	if !userExist {
		return userExist, userData, fmt.Errorf("登入驗證失敗")
	}
	return userExist, userData, nil
}

//token產生器
func generateToken(c *gin.Context, user Userinfo) {
	//產生SignKey: 簽名和解簽需要使用一個值
	j := mylib.NewJWT()

	//產生User claims資料
	claims := mylib.CustomClaims{
		UserAccount: user.Account,
		StandardClaims: jwtgo.StandardClaims{
			NotBefore: int64(time.Now().Unix() - 1000), //簽名生效時間
			ExpiresAt: int64(time.Now().Unix() + 3600), //簽名過期時間
			Issuer:    "Seaflower",                     //簽名頒發者
		},
	}

	//根據claims生成token
	token, err := j.CreateToken(claims)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  -1,
			"message": err.Error(),
			"data":    nil,
		})
	}

	log.Println(token)
	//取得User相關資料
	data := Logininfo{
		Account:  user.Account,
		Password: "",
		Token:    token,
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  0,
		"message": "登錄成功",
		"data":    data,
	})
	return
}
