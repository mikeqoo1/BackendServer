package router

import (
	mylib "BackendServer/mylib"
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

//用來判斷是哪種資料 0代表raw data 1代表form-data
var checkflag int

//Userinfo 管理員
type Userinfo struct {
	Account  string `json:"account"`
	Password string `json:"Password"`
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
