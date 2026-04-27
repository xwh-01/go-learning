package utils

import (
	"log"
	"os"
	"sync"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

var (
	logger      *log.Logger
	logLevel    LogLevel = INFO
	loggerMutex sync.RWMutex
	logFile     *os.File
)

func InitLogger(logFilePath string, level LogLevel) error {
	loggerMutex.Lock()
	defer loggerMutex.Unlock()

	if logFile != nil {
		logFile.Close()
	}

	var err error
	if logFilePath != "" {
		logFile, err = os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		logger = log.New(logFile, "", log.LstdFlags|log.Lshortfile)
	} else {
		logger = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	}

	logLevel = level
	return nil
}

func SetLogLevel(level LogLevel) {
	loggerMutex.Lock()
	defer loggerMutex.Unlock()
	logLevel = level
}

func GetLogLevel() LogLevel {
	loggerMutex.RLock()
	defer loggerMutex.RUnlock()
	return logLevel
}

func shouldLog(level LogLevel) bool {
	return level >= logLevel
}

func logMessage(level LogLevel, format string, v ...interface{}) {
	loggerMutex.RLock()
	defer loggerMutex.RUnlock()

	if !shouldLog(level) {
		return
	}

	prefix := ""
	switch level {
	case DEBUG:
		prefix = "[DEBUG] "
	case INFO:
		prefix = "[INFO] "
	case WARN:
		prefix = "[WARN] "
	case ERROR:
		prefix = "[ERROR] "
	case FATAL:
		prefix = "[FATAL] "
	}

	logger.Output(2, prefix+format)
}

func Debug(format string, v ...interface{}) {
	logMessage(DEBUG, format, v...)
}

func Info(format string, v ...interface{}) {
	logMessage(INFO, format, v...)
}

func Warn(format string, v ...interface{}) {
	logMessage(WARN, format, v...)
}

func LogError(format string, v ...interface{}) {
	logMessage(ERROR, format, v...)
}

func Fatal(format string, v ...interface{}) {
	logMessage(FATAL, format, v...)
	os.Exit(1)
}

func CloseLogger() error {
	loggerMutex.Lock()
	defer loggerMutex.Unlock()

	if logFile != nil {
		return logFile.Close()
	}
	return nil
}
