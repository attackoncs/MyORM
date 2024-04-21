package mylog

import (
	"io/ioutil"
	"log"
	"os"
	"sync"
)

// [info]颜色为蓝色，[error]为红色，log.Lshortfile支持显示文件名和代码行号
var (
	errorLog = log.New(os.Stdout, "\033[31m[error]\033[0m ", log.LstdFlags|log.Lshortfile)
	infoLog  = log.New(os.Stdout, "\033[34m[info ]\033[0m ", log.LstdFlags|log.Lshortfile)
	loggers  = []*log.Logger{errorLog, infoLog}
	mu       sync.Mutex
)

// log methods暴露四个方法
var (
	Error  = errorLog.Println
	Errorf = errorLog.Printf
	Info   = infoLog.Println
	Infof  = infoLog.Printf
)

// log level
const (
	InfoLevel = iota
	ErrorLevel
	Disabled
)

// 设置日志级别
func SetLevel(level int) {
	mu.Lock()
	defer mu.Unlock()

	for _, logger := range loggers {
		logger.SetOutput(os.Stdout)
	}

	if ErrorLevel < level {
		errorLog.SetOutput(ioutil.Discard) //不打印该日志
	}

	if InfoLevel < level {
		infoLog.SetOutput(ioutil.Discard) //不打印该日志
	}
}
