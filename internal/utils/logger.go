package utils

import (
	"io"
	"log"
	"os"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
)

// InitLogger 初始化日志记录器
func InitLogger(logFile io.Writer) {
	// 同时将日志输出到文件和控制台
	multiWriter := io.MultiWriter(os.Stdout, logFile)

	// 初始化信息日志记录器
	InfoLogger = log.New(multiWriter,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	// 初始化错误日志记录器
	ErrorLogger = log.New(multiWriter,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}
