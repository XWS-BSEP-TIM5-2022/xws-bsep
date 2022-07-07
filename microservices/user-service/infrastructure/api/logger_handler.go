package api

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/user_service/startup/config"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

type CustomLogger struct {
	WarningLogger *logrus.Entry
	InfoLogger    *logrus.Entry //*log.Logger //*zerolog.Logger
	ErrorLogger   *logrus.Entry
	SuccessLogger *logrus.Entry
	DebugLogger   *logrus.Entry
}

func NewCustomLogger() *CustomLogger {
	InfoLogger := setLogrusLogger(config.NewConfig().InfoLogsFile)
	ErrorLogger := setLogrusLogger(config.NewConfig().ErrorLogsFile)
	WarningLogger := setLogrusLogger(config.NewConfig().WarningLogsFile)
	SuccessLogger := setLogrusLogger(config.NewConfig().SuccessLogsFile)
	DebugLogger := setLogrusLogger(config.NewConfig().DebugLogsFile)

	return &CustomLogger{
		InfoLogger:    InfoLogger,
		ErrorLogger:   ErrorLogger,
		WarningLogger: WarningLogger,
		SuccessLogger: SuccessLogger,
		DebugLogger:   DebugLogger,
	}
}

func caller() func(*runtime.Frame) (function string, file string) {
	return func(f *runtime.Frame) (function string, file string) {
		p, _ := os.Getwd()
		return "", fmt.Sprintf("%s:%d", strings.TrimPrefix(f.File, p), f.Line)
	}
}

func setLogrusLogger(filename string) *logrus.Entry {
	mLog := logrus.New()
	mLog.SetReportCaller(true)
	// mLog.SetLevel(logrus.DebugLevel)

	logsFolderName := config.NewConfig().LogsFolder
	if _, err := os.Stat(logsFolderName); os.IsNotExist(err) {
		os.Mkdir(logsFolderName, 0777)
	}
	lumberjackLogger := &lumberjack.Logger{
		Filename:   logsFolderName + filename,
		MaxSize:    1,
		MaxBackups: 0, // datoteke nece biti obrisane
		MaxAge:     0,
		LocalTime:  true,
		Compress:   true,
	}
	mw := io.MultiWriter(os.Stdout, lumberjackLogger)
	mLog.SetOutput(mw)

	mLog.SetFormatter(&logrus.TextFormatter{ //TextFormatter //JSONFormatter
		CallerPrettyfier: caller(),
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyFile: "mehtod",
		},
		ForceColors: true,
	})
	contextLogger := mLog.WithFields(logrus.Fields{})
	return contextLogger
}

func (customLogger *CustomLogger) getFileSize() {
	logsFolderName := config.NewConfig().LogsFolder
	filename := config.NewConfig().InfoLogsFile
	file, err := os.OpenFile(logsFolderName+filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	fi, err := file.Stat()
	if err != nil {
		log.Panic("error")
	}
	fmt.Println("----------------------------------------")
	fmt.Println("size: ", fi.Size())
	fmt.Println("----------------------------------------")
}
