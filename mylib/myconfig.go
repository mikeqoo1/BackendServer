package mylib

import (
	viper "github.com/spf13/viper"
)

//myConfigStruct 我自定義的config結構
type myConfigStruct struct {
	TimeZone string
	Interval int
	Host     string
	Port     int
	RpcHost  string
	RpcPort  int
	LogFile  string
	State    string

	//天坦api
	TendasoftHost string
	TendasoftPort int

	//GRPC
	GrpcHost string
	GrpcPort int

	//資料庫
	DBHost         string
	DBPort         int
	DBUser         string
	DBPassWord     string
	DBDataBase     string
	DBMaxOpenConns int
	DBMaxIdleConns int

	//Https憑證
	Cert    string
	Privkey string

	//開發商利潤
	Profit float64
}

//MyConfig 自定義的config變數
var MyConfig myConfigStruct

//UsaTimeMap 美國時區地圖
var UsaTimeMap map[string]string

//InitMyConfig 我的設定檔
func InitMyConfig(configpath string) {
	viper.SetConfigName("config") //指定文件的名稱
	viper.AddConfigPath(configpath)
	err := viper.ReadInConfig() //根據以上定讀取文件
	if err != nil {
		MyLogger.Fatal("Fatal error config file" + err.Error())
	}

	MyConfig.TimeZone = viper.GetString("BASIC.timeZone")
	MyConfig.Interval = viper.GetInt("BASIC.interval")
	MyConfig.Host = viper.GetString("BASIC.host")
	MyConfig.Port = viper.GetInt("BASIC.port")
	MyConfig.RpcHost = viper.GetString("BASIC.rpchost")
	MyConfig.RpcPort = viper.GetInt("BASIC.rpcport")
	MyConfig.LogFile = viper.GetString("BASIC.logfile")
	MyConfig.State = viper.GetString("BASIC.state")

	MyConfig.TendasoftHost = viper.GetString("Tendasoft.host")
	MyConfig.TendasoftPort = viper.GetInt("Tendasoft.port")

	MyConfig.GrpcHost = viper.GetString("GRPC.host")
	MyConfig.GrpcPort = viper.GetInt("GRPC.port")

	MyConfig.DBHost = viper.GetString("DB.host")
	MyConfig.DBPort = viper.GetInt("DB.port")
	MyConfig.DBUser = viper.GetString("DB.user")
	MyConfig.DBPassWord = viper.GetString("DB.password")
	MyConfig.DBDataBase = viper.GetString("DB.database")
	MyConfig.DBMaxOpenConns = viper.GetInt("DB.MaxOpenConns")
	MyConfig.DBMaxIdleConns = viper.GetInt("DB.MaxIdleConns")

	MyConfig.Cert = viper.GetString("BASIC.cert")
	MyConfig.Privkey = viper.GetString("BASIC.privkey")

	MyConfig.Profit = viper.GetFloat64("JoyGame.profit")
}

//InitUsaTimeZoneMap 初始美國地區map
func InitUsaTimeZoneMap() {
	UsaTimeMap = make(map[string]string)
	UsaTimeMap["Connecticut"] = "UTC-5"
	UsaTimeMap["Delaware"] = "UTC-5"
	UsaTimeMap["District of Columbia"] = "UTC-5"
	UsaTimeMap["Florida"] = "UTC-5"
	UsaTimeMap["Georgia"] = "UTC-5"
	UsaTimeMap["Indiana"] = "UTC-5"
	UsaTimeMap["Kentucky"] = "UTC-5"
	UsaTimeMap["Maine"] = "UTC-5"
	UsaTimeMap["Maryland"] = "UTC-5"
	UsaTimeMap["Massachusetts"] = "UTC-5"
	UsaTimeMap["Michigan"] = "UTC-5"
	UsaTimeMap["New Hampshire"] = "UTC-5"
	UsaTimeMap["New Jersey"] = "UTC-5"
	UsaTimeMap["New York state"] = "UTC-5"
	UsaTimeMap["North Carolina"] = "UTC-5"
	UsaTimeMap["Ohio"] = "UTC-5"
	UsaTimeMap["Oklahoma"] = "UTC-5"
	UsaTimeMap["Pennsylvania"] = "UTC-5"
	UsaTimeMap["Rhode Island"] = "UTC-5"
	UsaTimeMap["South Carolina"] = "UTC-5"
	UsaTimeMap["Vermont"] = "UTC-5"
	UsaTimeMap["Virginia"] = "UTC-5"
	UsaTimeMap["West Virginia"] = "UTC-5"
	UsaTimeMap["Alabama"] = "UTC-6"
	UsaTimeMap["Arkansas"] = "UTC-6"
	UsaTimeMap["Idaho"] = "UTC-6"
	UsaTimeMap["Illinois"] = "UTC-6"
	UsaTimeMap["Iowa"] = "UTC-6"
	UsaTimeMap["Kansas"] = "UTC-6"
	UsaTimeMap["Louisiana"] = "UTC-6"
	UsaTimeMap["Minnesota"] = "UTC-6"
	UsaTimeMap["Mississippi"] = "UTC-6"
	UsaTimeMap["Missouri"] = "UTC-6"
	UsaTimeMap["Nebraska"] = "UTC-6"
	UsaTimeMap["North Dakota"] = "UTC-6"
	UsaTimeMap["South Dakota"] = "UTC-6"
	UsaTimeMap["Tennessee"] = "UTC-6"
	UsaTimeMap["Texas"] = "UTC-6"
	UsaTimeMap["Wisconsin"] = "UTC-6"
	UsaTimeMap["Arizona"] = "UTC-7"
	UsaTimeMap["Colorado"] = "UTC-7"
	UsaTimeMap["Montana"] = "UTC-7"
	UsaTimeMap["New Mexico"] = "UTC-7"
	UsaTimeMap["Utah"] = "UTC-7"
	UsaTimeMap["Wyoming"] = "UTC-7"
	UsaTimeMap["California"] = "UTC-8"
	UsaTimeMap["Nevada"] = "UTC-8"
	UsaTimeMap["Oregon"] = "UTC-8"
	UsaTimeMap["Washington state"] = "UTC-8"
	UsaTimeMap["Alaska"] = "UTC-9"
	UsaTimeMap["Hawaii"] = "UTC-10"
}
