package main

import (
	mylib "BackendServer/mylib"
	url "BackendServer/router"
	"fmt"

	"github.com/gin-gonic/gin"
)

//常用變數
var err error

func init() {
	mylib.InitMyConfig()
	mylib.MyLog()
	mylib.InitDBpool()
}

func main() {
	addr := fmt.Sprintf("%s:%d", "127.0.0.1", mylib.MyConfig.ServerPort)
	router := gin.Default()

	v1 := router.Group("/test")
	{
		v1.GET("/login", url.Login)
		// v1.POST("/submit", submitEndpoint)
		// v1.POST("/read", readEndpoint)
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
