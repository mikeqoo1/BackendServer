package router

import (
	mylib "BackendServer/mylib"
	"bytes"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

//用來判斷是哪種資料 1代表raw data 0代表form-data
var checkflag int

//Admin 管理員
type Admin struct {
	Account string `json:"account"`
	Status  string `json:"status"`
}

// //Routerurl 存放URL的設定結構
// type Routerurl struct {
// 	URL路徑 string
// 資料結構
// }

// //NewRouterURL new a Routerurl object
// func NewRouterURL() *Routerurl {
// 	routerURL := new(Routerurl)
// 	return routerURL
// }

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

//API api格式
func API(c *gin.Context) {
	c.JSON(http.StatusOK, "hello word")
}

//API1 api格式
func API1(c *gin.Context) {
	result, err := mylib.ConstDBpool.Exec("SQL語法")
	if err != nil {
		mylib.MyLogger.Error("API1 SQL 失敗" + err.Error())
	}
	mylib.MyLogger.Debug("result:", result)
	mylib.MyLogger.Debug("user IP:", c.ClientIP())

	c.JSON(http.StatusOK, "hello word")
}

//API2 api格式
func API2(c *gin.Context) {
	mylib.MyLogger.Debug("誰來新增: ", c.ClientIP())
	stmt, err := mylib.ConstDBpool.Prepare("INSERT INTO Admin SET Account = ?, createtime = ?, status = ?")
	if err != nil {
		mylib.MyLogger.Error("API1 SQL 失敗", err.Error())
	}

	var info Admin

	//先用form-data解析
	user := c.DefaultPostForm("account", "NULL")
	info.Account = user
	if user == "NULL" {
		mylib.MyLogger.Warn("該筆訊息不是form-data")
		checkflag = 1
	} else {
		status := c.DefaultPostForm("status", "0")
		info.Status = status
	}

	if checkflag == 1 {
		err = c.BindJSON(&info)
		if err != nil {
			mylib.MyLogger.Error("API1 BindJSON 失敗 ", err.Error())
		}
		mylib.MyLogger.Debug("info: ", info)
	}

	_, err1 := stmt.Exec(info.Account, time.Now(), info.Status)
	if err1 != nil {
		mylib.MyLogger.Error("New Error:", err1.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err1.Error(),
		})
	} else {
		mylib.MyLogger.Info("New success")
		c.JSON(http.StatusOK, gin.H{
			"message": "Create success",
		})
	}

}
