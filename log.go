package xin

import (
	"fmt"
	"io"
	"log"
	"os"
)

// Logger 定义了日志接口
type Logger interface {
	Debug(v ...any)
	Debugf(format string, v ...any)
	Info(v ...any)
	Infof(format string, v ...any)
	Error(v ...any)
	Errorf(format string, v ...any)
}

// stdLogger 是标准输出的日志实现
type stdLogger struct {
	infoLog  *log.Logger
	debugLog *log.Logger
	errorLog *log.Logger
	debug    bool // 是否启用调试日志
}

// NewStdLogger 创建一个标准输出的日志实现
func NewStdLogger() Logger {
	return &stdLogger{
		infoLog:  log.New(os.Stdout, "[INFO] ", log.LstdFlags),
		debugLog: log.New(os.Stdout, "[DEBUG] ", log.LstdFlags),
		errorLog: log.New(os.Stderr, "[ERROR] ", log.LstdFlags|log.Lshortfile),
		debug:    false,
	}
}

// NewCustomLogger 创建一个自定义输出的日志实现
func NewCustomLogger(infoOut, debugOut, errorOut io.Writer) Logger {
	return &stdLogger{
		infoLog:  log.New(infoOut, "[INFO] ", log.LstdFlags),
		debugLog: log.New(debugOut, "[DEBUG] ", log.LstdFlags),
		errorLog: log.New(errorOut, "[ERROR] ", log.LstdFlags|log.Lshortfile),
		debug:    false,
	}
}

func (l *stdLogger) Info(v ...any) {
	l.infoLog.Output(2, fmt.Sprint(v...))
}

func (l *stdLogger) Infof(format string, v ...any) {
	l.infoLog.Output(2, fmt.Sprintf(format, v...))
}

func (l *stdLogger) Debug(v ...any) {
	if !l.debug {
		return
	}
	l.debugLog.Output(2, fmt.Sprint(v...))
}

func (l *stdLogger) Debugf(format string, v ...any) {
	if !l.debug {
		return
	}
	l.debugLog.Output(2, fmt.Sprintf(format, v...))
}

func (l *stdLogger) Error(v ...any) {
	l.errorLog.Output(2, fmt.Sprint(v...))
}

func (l *stdLogger) Errorf(format string, v ...any) {
	l.errorLog.Output(2, fmt.Sprintf(format, v...))
}

// SetDebug 设置是否启用调试日志
func (l *stdLogger) SetDebug(debug bool) {
	l.debug = debug
}

// 默认的日志实例
var defaultLogger = NewStdLogger()

// SetLogger 设置全局默认的日志实例
func SetLogger(logger Logger) {
	defaultLogger = logger
}

// GetLogger 获取全局默认的日志实例
func GetLogger() Logger {
	return defaultLogger
}

func LogInfo(v ...any) {
	defaultLogger.Info(v...)
}

func LogInfof(format string, v ...any) {
	defaultLogger.Infof(format, v...)
}

func LogDebug(v ...any) {
	defaultLogger.Debug(v...)
}

func LogDebugf(format string, v ...any) {
	defaultLogger.Debugf(format, v...)
}

func LogError(v ...any) {
	defaultLogger.Error(v...)
}

func LogErrorf(format string, v ...any) {
	defaultLogger.Errorf(format, v...)
}
