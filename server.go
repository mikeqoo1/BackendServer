package main

import (
	mylib "BackendServer/mylib"
	url "BackendServer/router"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

//常用變數
var err error

func init() {
	mylib.InitMyConfig()
	mylib.MyLog()
	mylib.InitDBpool()
}

//crosHandler 處理跨域問題
func crosHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin") //請求頭部
		if origin != "" {
			//接收客戶端傳送的origin (重要)
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			//伺服器支援的所有跨域請求的方法
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			//允許跨域設定可以返回其他子段，可以自定義欄位
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session, "+
				"X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, "+
				"X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, "+
				"Content-Type, Pragma, token, openid, opentoken")
			//允許瀏覽器(客戶端)可以解析的頭部 (重要)
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, "+
				"Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type, "+
				"Expires, Last-Modified, Pragma, FooBar")
			//設定快取時間
			c.Header("Access-Control-Max-Age", "172800")
			//允許客戶端傳遞校驗資訊比如 cookie (重要)
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		//允許型別校驗
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "ok!")
		}

		c.Next()
	}
}

func main() {
	addr := fmt.Sprintf("%s:%d", "127.0.0.1", mylib.MyConfig.ServerPort)
	router := gin.Default()
	router.Use(crosHandler())
	router.Use(url.CheckData())

	v1 := router.Group("/apis")
	{
		v1.POST("/register", url.Registeruser)
		v1.POST("/login", url.Login)
	}

	t := router.Group("/apis/test")
	{
		t.GET("/user", url.List)       //列出所有user
		t.POST("/user", url.New)       //新建一個user
		t.PUT("/user/ID", url.Update)  //更新某個指定user的信息（提供該user的全部信息）
		t.DELETE("/user/ID", url.Drop) //刪除某個user
	}

	// v2 := router.Group("/v2")
	// {
	// 	v2.POST("/login", loginEndpoint)
	// 	v2.POST("/submit", submitEndpoint)
	// 	v2.POST("/read", readEndpoint)
	// }

	err = router.Run(addr)
	if err != nil {
		mylib.MyLogger.Error("BackendServer啟動失敗" + err.Error())
	}
}
