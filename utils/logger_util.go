package utils

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// LogToFile saves the request details into a file with year/month/day folder structure
func LogToFile(prefix string, headers http.Header, body []byte) {
	now := time.Now()

	// Format: logs/2026/02/19
	dirPath := filepath.Join("logs", now.Format("2006"), now.Format("01"), now.Format("02"))

	// Create directories if they don't exist
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		log.Printf("Error creating log directory: %v", err)
		return
	}

	fileName := fmt.Sprintf("%s.log", prefix)
	filePath := filepath.Join(dirPath, fileName)

	// Open file in append mode, create if not exists
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error opening log file: %v", err)
		return
	}
	defer f.Close()

	// Prepare log entry
	logEntry := fmt.Sprintf("--- [%s] ---\n", now.Format("15:04:05"))
	logEntry += "Headers:\n"
	for k, v := range headers {
		logEntry += fmt.Sprintf("  %s: %v\n", k, v)
	}
	logEntry += fmt.Sprintf("Body: %s\n\n", string(body))

	if _, err := f.WriteString(logEntry); err != nil {
		log.Printf("Error writing to log file: %v", err)
	}
}

// LogEvent saves a simple string message into a file with year/month/day folder structure
func LogEvent(prefix string, message string) {
	now := time.Now()
	dirPath := filepath.Join("logs", now.Format("2006"), now.Format("01"), now.Format("02"))

	if err := os.MkdirAll(dirPath, 0755); err != nil {
		log.Printf("Error creating log directory: %v", err)
		return
	}

	filePath := filepath.Join(dirPath, fmt.Sprintf("%s.log", prefix))
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error opening log file: %v", err)
		return
	}
	defer f.Close()

	logEntry := fmt.Sprintf("[%s] %s\n", now.Format("15:04:05"), message)
	if _, err := f.WriteString(logEntry); err != nil {
		log.Printf("Error writing to log file: %v", err)
	}
}
