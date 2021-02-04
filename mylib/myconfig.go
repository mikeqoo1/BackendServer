package mylib

import (
	viper "github.com/spf13/viper"
)

//myConfigStruct 我自定義的config結構
type myConfigStruct struct {
	TomeZone         string
	Interval         int
	ServerPort       int

	DBHost         string
	DBPort         int
	DBUser         string
	DBPassWord     string
	DBDataBase     string
	DBMaxOpenConns int
	DBMaxIdleConns int
}

//MyConfig 自定義的config變數
var MyConfig myConfigStruct

//InitMyConfig 我的設定檔
func InitMyConfig() {
	viper.SetConfigName("config")                         // 指定文件的名稱
	viper.AddConfigPath("/Projects/BackendServer/config") // 配置文件和執行檔目錄
	err := viper.ReadInConfig()                           // 根據以上定讀取文件
	if err != nil {
		MyLogger.Fatal("Fatal error config file" + err.Error())
	}
	MyConfig.TomeZone = viper.GetString("BASIC.timeZone")
	MyConfig.Interval = viper.GetInt("BASIC.interval")
	MyConfig.ServerPort = viper.GetInt("BASIC.serverPort")

	MyConfig.DBHost = viper.GetString("DB.host")
	MyConfig.DBPort = viper.GetInt("DB.port")
	MyConfig.DBUser = viper.GetString("DB.user")
	MyConfig.DBPassWord = viper.GetString("DB.password")
	MyConfig.DBDataBase = viper.GetString("DB.database")
	MyConfig.DBMaxOpenConns = viper.GetInt("DB.MaxOpenConns")
	MyConfig.DBMaxIdleConns = viper.GetInt("DB.MaxIdleConns")
}
