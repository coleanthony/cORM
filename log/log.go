package log

import (
	"io"
	"log"
	"os"
	"sync"
)

/*
支持日志分级（Info、Error、Disabled 三级）。
不同层级日志显示时使用不同的颜色区分。
显示打印日志代码对应的文件名和行号。
*/
//[info ] 颜色为蓝色，[error] 为红色。
//使用 log.Lshortfile 支持显示文件名和代码行号。
var (
	errorLog = log.New(os.Stdout, "\033[31m[error]\033[0m ", log.LstdFlags|log.Lshortfile)
	infoLog  = log.New(os.Stdout, "\033[34m[info ]\033[0m ", log.LstdFlags|log.Lshortfile)
	loggers  = []*log.Logger{errorLog, infoLog}
	mutex    sync.Mutex
)

// log methods
// 暴露 Error，Errorf，Info，Infof 4个方法
var (
	Error  = errorLog.Println
	Errorf = errorLog.Printf
	Info   = infoLog.Println
	Infof  = infoLog.Printf
)

// log levels
const (
	InfoLevel = iota
	ErrorLevel
	Disabled
)

// 如果设置为 ErrorLevel，infoLog 的输出会被定向到 ioutil.Discard，即不打印该日志
func SetLevel(level int) {
	mutex.Lock()
	defer mutex.Unlock()
	for _, logger := range loggers {
		logger.SetOutput(os.Stdout)
	}
	if ErrorLevel < level {
		errorLog.SetOutput(io.Discard)
	}
	if InfoLevel < level {
		infoLog.SetOutput(io.Discard)
	}
}
