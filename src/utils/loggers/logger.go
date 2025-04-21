package loggers

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

type Logger struct {
	file *os.File
	log  *log.Logger
}

func NewLogger() *Logger {
	// Buat folder logs jika belum ada
	logDir := "storage/logs"
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		os.Mkdir(logDir, os.ModePerm)
	}

	// Gunakan nama file berdasarkan tanggal
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

func (l *Logger) Info(message string, context map[string]interface{}) {
	l.writeLog("INFO", message, context)
}

func (l *Logger) Error(message string, context map[string]interface{}) {
	l.writeLog("ERROR", message, context)
}

func (l *Logger) writeLog(level, message string, context map[string]interface{}) {
	timestamp := time.Now().Format(time.RFC3339)
	entry := map[string]interface{}{
		"context":   context,
		"message":   message,
		"level":     level,
		"timestamp": timestamp,
	}

	jsonData, err := json.Marshal(entry)
	if err != nil {
		fmt.Printf("Failed to marshal log entry: %v", err)
		return
	}

	l.log.Println(string(jsonData))
}

func (l *Logger) Close() {
	l.file.Close()
}
