package loggers

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	Log  *Logger
	once sync.Once
)

type Logger struct {
	file *os.File
	log  *log.Logger
}

// init dijalankan otomatis saat package ini di-import
func init() {
	once.Do(func() {
		Log = NewLogger()
	})
}

// NewLogger initializes and returns a new Logger instance. It ensures that
// the log directory exists, creates or opens a log file with the current date
// as its name, and sets up the file for appending log entries. If the log file
// cannot be opened, the function will terminate the program.

func NewLogger() *Logger {
	logDir := "src/storage/logs"

	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		_ = os.MkdirAll(logDir, os.ModePerm)
	}

	fileName := time.Now().Format("2006-01-02") + ".log"
	filePath := filepath.Join(logDir, fileName)

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	return &Logger{
		file: file,
		log:  log.New(file, "", 0),
	}
}

// Info logs the given message at the INFO level with the provided context.
// The context is stored as structured data in the log entry.
func (l *Logger) Info(message string, context map[string]interface{}) {
	l.writeLog("INFO", message, context)
}

// Error logs the given message at the ERROR level with the provided context.
// The context is stored as structured data in the log entry.

func (l *Logger) Error(message string, context map[string]interface{}) {
	l.writeLog("ERROR", message, context)
}

// writeLog logs the given message with the provided context and level.
// It marshals the log entry into a JSON string and writes it to the log file,
// surrounded by a separator. If marshaling fails, it prints an error message
// and returns.
func (l *Logger) writeLog(level, message string, context map[string]interface{}) {
	timestamp := time.Now().Format(time.RFC3339)
	entry := map[string]interface{}{
		"timestamp": timestamp,
		"level":     level,
		"message":   message,
		"context":   context,
	}

	jsonData, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		fmt.Printf("Failed to marshal log entry: %v", err)
		return
	}

	separator := "\n============================================================\n"

	logOutput := fmt.Sprintf("%s\n%s%s", separator, string(jsonData), separator)

	l.log.Println(logOutput)
}

func (l *Logger) Close() {
	l.file.Close()
}
