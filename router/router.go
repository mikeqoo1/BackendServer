package router

import (
	mylib "BackendServer/mylib"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

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

//Login Login Test
func Login(c *gin.Context) {
	result, err := mylib.ConstDBpool.Exec(
		"INSERT INTO bet_cluster (ServerID, PlatformID, MemberCode, AgentID, LobbyID, GameID, UserID, ThirdPartyUserID," +
			"ThirdPartyUserIDStr, Account, Agent, Currency, OrderID, Round, Bet, Win, WinLose, StartTime, EndTime, OrderState, IsProcess)" +
			"VALUES(0, 0, 0, 0, 0, 0, 0, 0, 'F', 'A', 'B', 'RD', 'C', 0, 0, 0, 0, '2020-10-20', '2020-10-20', 0, 0)")
	if err != nil {
		mylib.MyLogger.Error("Login SQL 失敗" + err.Error())
		fmt.Println(err.Error())
	}
	fmt.Println(result)

	//fmt.Println(c.Params)

	c.JSON(http.StatusOK, "hello word")
}
