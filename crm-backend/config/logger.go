package config

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type Logger struct {
	file *os.File
}

var (
	loggerInstance *Logger
	once           sync.Once
)

func GetLoggerInstance() *Logger {
	once.Do(func() {
		file, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatalf("Failed to open log file: %v", err)
		}
		loggerInstance = &Logger{file: file}
	})
	return loggerInstance
}

func (l *Logger) Log(message, errMsg, code string) {
	logMessage := fmt.Sprintf("message: %s%s code: %s, at %v", message, errMsg, code, time.Now())
	_, err := l.file.WriteString(logMessage + "\n")
	if err != nil {
		log.Fatalf("Failed to write to log file: %v", err)
	}
}
