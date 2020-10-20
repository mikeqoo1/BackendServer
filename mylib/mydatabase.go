package mylib

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var err error

//ConstDBpool 共用資料庫連線
var ConstDBpool *sql.DB

//DataRow ...
type DataRow map[string]interface{}

//RowCollection 把取得的資料用map存起來
type RowCollection []DataRow

//DataColumn 欄位名稱
type DataColumn []string

//DataTable 結果的結構
type DataTable struct {
	Rows    RowCollection
	Columns DataColumn
}

//FromRows 把撈出來的資料集轉換成struct的格式存取
func FromRows(rows *sql.Rows) (DataTable, error) {

	dt := DataTable{}

	columns, err := rows.Columns()
	if err != nil {
		return dt, err
	}

	dt.Columns = columns
	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		dt.Rows = append(dt.Rows, entry)
	}
	return dt, nil
}

//InitDBpool 建立一個共用的DB物件給大家用
func InitDBpool() {
	var connectionString = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?allowNativePasswords=true",
		MyConfig.DBUser, MyConfig.DBPassWord, MyConfig.DBHost, MyConfig.DBPort, MyConfig.DBDataBase)
	ConstDBpool, err = sql.Open("mysql", connectionString)
	if err != nil {
		MyLogger.Error("SQL open 失敗" + err.Error())
	}
	fmt.Println("DB config", MyConfig.DBMaxOpenConns, MyConfig.DBMaxIdleConns)
	ConstDBpool.SetMaxOpenConns(MyConfig.DBMaxOpenConns)
	ConstDBpool.SetMaxIdleConns(MyConfig.DBMaxIdleConns)
	err = ConstDBpool.Ping()
	if err != nil {
		MyLogger.Error("ConstDBpool.Ping Error" + err.Error())
	}
}
