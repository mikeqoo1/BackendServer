package main

import (
	mylib "BackendServer/mylib"
	"fmt"
	"net/http"
)

func init() {
	mylib.InitMyConfig()
	mylib.MyLog()
}

func main() {
	addr := fmt.Sprintf("%s:%d", "127.0.0.1", mylib.MyConfig.ServerPort)

	mylib.MyLogger.Debug("正在listen" + addr)

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		mylib.MyLogger.Fatal("監聽失敗" + err.Error())
	}
}
