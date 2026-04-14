package logger

import (
	"log"
	"os"
	"sync"
)

type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

var (
	currentLevel LogLevel = LevelInfo
	logger       *log.Logger
	mu           sync.RWMutex
)

// Init 初始化日志
func Init(debug bool) error {
	mu.Lock()
	defer mu.Unlock()

	if debug {
		currentLevel = LevelDebug
	} else {
		currentLevel = LevelInfo
	}

	// 创建日志文件
	logFile, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		// 如果创建文件失败，使用标准输出
		logger = log.New(os.Stdout, "[WeChatPadPro] ", log.LstdFlags|log.Lshortfile)
		return err
	}

	logger = log.New(logFile, "[WeChatPadPro] ", log.LstdFlags|log.Lshortfile)
	return nil
}

// Debug 记录调试日志
func Debug(format string, args ...interface{}) {
	mu.RLock()
	defer mu.RUnlock()

	if currentLevel <= LevelDebug && logger != nil {
		logger.Printf("[DEBUG] "+format, args...)
	}
}

// Info 记录信息日志
func Info(format string, args ...interface{}) {
	mu.RLock()
	defer mu.RUnlock()

	if currentLevel <= LevelInfo && logger != nil {
		logger.Printf("[INFO] "+format, args...)
	}
}

// Warn 记录警告日志
func Warn(format string, args ...interface{}) {
	mu.RLock()
	defer mu.RUnlock()

	if currentLevel <= LevelWarn && logger != nil {
		logger.Printf("[WARN] "+format, args...)
	}
}

// Error 记录错误日志
func Error(format string, args ...interface{}) {
	mu.RLock()
	defer mu.RUnlock()

	if currentLevel <= LevelError && logger != nil {
		logger.Printf("[ERROR] "+format, args...)
	}
}

// Fatal 记录致命日志并退出
func Fatal(format string, args ...interface{}) {
	mu.RLock()
	defer mu.RUnlock()

	if currentLevel <= LevelFatal && logger != nil {
		logger.Fatalf("[FATAL] "+format, args...)
	}
}

// SetLevel 设置日志级别
func SetLevel(level LogLevel) {
	mu.Lock()
	defer mu.Unlock()
	currentLevel = level
}