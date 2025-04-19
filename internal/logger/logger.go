package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Logger структура для логування
type Logger struct {
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
	DebugLogger *log.Logger
}

// NewLogger створює новий логер з заданими налаштуваннями
func NewLogger() *Logger {
	infoFile, err := os.OpenFile("logs/info.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		// Якщо не можемо створити файл, використовуємо stdout
		infoFile = os.Stdout
	}

	errorFile, err := os.OpenFile("logs/error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		// Якщо не можемо створити файл, використовуємо stderr
		errorFile = os.Stderr
	}

	// Створюємо логери для різних рівнів логування
	infoLogger := log.New(infoFile, "INFO: ", log.Ldate|log.Ltime)
	errorLogger := log.New(errorFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	debugLogger := log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)

	return &Logger{
		InfoLogger:  infoLogger,
		ErrorLogger: errorLogger,
		DebugLogger: debugLogger,
	}
}

// Info логує інформаційне повідомлення
func (l *Logger) Info(message string, args ...interface{}) {
	timestamp := time.Now().Format(time.RFC3339)
	if len(args) > 0 {
		l.InfoLogger.Printf("[%s] %s: %v", timestamp, message, args)
	} else {
		l.InfoLogger.Printf("[%s] %s", timestamp, message)
	}
}

// Error логує повідомлення про помилку
func (l *Logger) Error(message string, err error, args ...interface{}) {
	timestamp := time.Now().Format(time.RFC3339)
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}

	if len(args) > 0 {
		argsStr := fmt.Sprint(args...)
		l.ErrorLogger.Printf("[%s] %s: %s, args: %s", timestamp, message, errMsg, argsStr)
	} else {
		l.ErrorLogger.Printf("[%s] %s: %s", timestamp, message, errMsg)
	}
}

// Debug логує повідомлення для відлагодження
func (l *Logger) Debug(message string, args ...interface{}) {
	timestamp := time.Now().Format(time.RFC3339)
	if len(args) > 0 {
		l.DebugLogger.Printf("[%s] %s: %v", timestamp, message, args)
	} else {
		l.DebugLogger.Printf("[%s] %s", timestamp, message)
	}
}
