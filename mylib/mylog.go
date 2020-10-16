package mylib

import (
	"fmt"
	"os"
	"strconv"
	"time"

	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"

	"path"

	logrus "github.com/sirupsen/logrus"
)

//MyLogger 我自定義的log變數
var MyLogger *logrus.Entry

//MyLog 我的log初始化
func MyLog() {
	log := logrus.New()

	dateStr := "/Projects/BackendServer/log/BackendServer" + time.Now().Format("2006-01-02") + ".log"

	logFile, err := os.OpenFile(dateStr, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error 開啟 log 失敗", err.Error())
	}

	log.AddHook(newRotateHook("", dateStr, 31*24*time.Hour, time.Hour))

	// 為當前logrus例項設定訊息的輸出, 可以設定logrus例項的輸出到任意io.writer EX: os.Stdout
	log.Out = logFile

	//為當前logrus例項設定訊息輸出格式為json格式
	formatter := &logrus.JSONFormatter{
		// time格式
		TimestampFormat: time.StampNano,
	}

	// // formatter := &logrus.TextFormatter{
	// // 	// time格式
	// // 	TimestampFormat: time.StampNano,
	// // }

	log.SetFormatter(formatter)

	//檔出終端機
	log.SetOutput(os.Stdout)

	//設定log等級
	log.SetLevel(logrus.DebugLevel)

	//設定log欄位
	MyLogger = log.WithFields(logrus.Fields{
		"LogTag": "LogValue",
		"Key":    "Value",
		"Key1":   "Value1",
	})

}

func newRotateHook(logPath string, logFileName string, maxAge time.Duration, rotationTime time.Duration) *lfshook.LfsHook {
	baseLogPath := path.Join(logPath, logFileName)

	writer, err := rotatelogs.New(
		baseLogPath+strconv.Itoa(time.Now().Hour()), // 檔名後面加上時
		rotatelogs.WithLinkName(baseLogPath),        // 生成軟連結, 指向最新log
		rotatelogs.WithMaxAge(maxAge),               // log最大保存時間
		rotatelogs.WithRotationTime(rotationTime),   // log切割的時間間隔
	)
	if err != nil {
		MyLogger.Errorf("config local file system logger error. %+v", errors.WithStack(err))
	}
	return lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer, // 設定不同等級的輸出目的
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, &logrus.JSONFormatter{TimestampFormat: time.StampNano}) //初始logrus訊息輸出格式為json格式
}
