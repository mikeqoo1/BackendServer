package main

import (
	md "BackendServer/middleware"
	mylib "BackendServer/mylib"
	url "BackendServer/router"

	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
)

//常用變數
var (
	//IsHttps https版本
	IsHttps string
	//IsJC JC版本
	IsJC string
	//IsDebug debug訊息
	IsDebug    string
	configpath string
)

func init() {
	if IsJC == "TURE" {
		configpath = "/data/joygame/config"
		gin.SetMode(gin.ReleaseMode)
	} else {
		configpath = "/joygame/backendserver/config"
	}
	mylib.MyLog()
	mylib.InitDBpool()
}

//TLSHandler http to https
func TLSHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		secureMiddleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     "0.0.0.0:8080",
		})
		err := secureMiddleware.Process(c.Writer, c.Request)

		// If there was an error, do not continue.
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		c.Next()
	}
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
	var err error
	addr := fmt.Sprintf("%s:%d", mylib.MyConfig.Host, mylib.MyConfig.Port)
	router := gin.Default()
	router.Use(crosHandler())
	if IsJC == "TURE" {
		router.Use(TLSHandler())
	}
	router.Use(url.CheckData())

	router.GET("/develop/check", url.CheckToken)

	backend := router.Group("/")
	backend.Use(md.JWTAuth())
	{

	}

	pos := router.Group("/pos")
	pos.Use(md.JWTAuth())
	{
		v1.POST("/register", url.Registeruser)
		v1.POST("/login", url.Login)
	}

	group := router.Group("/group")
	group.Use(md.JWTAuth())
	{
		group.GET("/", url.List)    //列出grouplist
		group.POST("/", url.New)    //新建group
		group.PUT("/", url.Update)  //更新group
		group.DELETE("/", url.Drop) //刪除group

		group.POST("/dd", url.New)   //新建IGS版本的分潤
		group.PUT("/dd", url.Update) //更新IGS版本的分潤

	}

	store := router.Group("/store")
	store.Use(md.JWTAuth())
	{
		store.GET("/", url.List)    //列出store
		store.POST("/", url.New)    //新建store
		store.PUT("/", url.Update)  //更新store
		store.DELETE("/", url.Drop) //刪除store
	}

	t := router.Group("/apis/test")
	{
		t.GET("/user", url.List)       //列出所有user
		t.POST("/user", url.New)       //新建一個user
		t.PUT("/user/ID", url.Update)  //更新某個指定user的信息（提供該user的全部信息）
		t.DELETE("/user/ID", url.Drop) //刪除某個user
	}

	err = router.Run(addr)
	if err != nil {
		mylib.MyLogger.Error("BackendServer啟動失敗" + err.Error())
	}
}
