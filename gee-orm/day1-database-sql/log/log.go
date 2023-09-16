package log

import (
	"io"
	"log"
	"os"
	"sync"
)

var (
	// 红色
	errorLog = log.New(os.Stdout, "\033[31m[error]\033[0m ", log.LstdFlags|log.Lshortfile)
	// 蓝色
	infoLog = log.New(os.Stdout, "\033[34m[info ]\033[0m ", log.LstdFlags|log.Lshortfile)
	loggers = []*log.Logger{errorLog, infoLog}
	// mu 是一个用于多线程同步的互斥锁（Mutex），用于在多线程环境中保护共享的资源，这里主要是保护 loggers 切片
	mu sync.Mutex
)

// log methods
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

// SetLevel 控制日志级别，可以决定应用程序在不同的情况下记录哪些级别的日志
func SetLevel(level int) {
	mu.Lock()
	defer mu.Unlock()

	for _, logger := range loggers {
		logger.SetOutput(os.Stdout)
	}

	if ErrorLevel < level {
		// 如果 ErrorLevel 小于 level，则将 errorLog 的输出设置为 ioutil.Discard，这意味着所有的错误级别日志消息将被丢弃，不会输出到终端
		errorLog.SetOutput(io.Discard)
	}
	if InfoLevel < level {
		infoLog.SetOutput(io.Discard)
	}
}
