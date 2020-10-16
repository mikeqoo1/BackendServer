#主程式
Main=server.go

#執行檔
Server_Out=-o bin/server.out

main:
	go install
	go build $(Server_Out) $(Main)

