package logger

import (
	"log"
	"os"
	"time"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
)

func init() {
	// Create logs directory if it doesn't exist
	if err := os.MkdirAll("logs", 0755); err != nil {
		log.Fatal("Failed to create logs directory:", err)
	}

	// Open log file with current date
	currentTime := time.Now()
	logFile, err := os.OpenFile(
		"logs/"+currentTime.Format("2006-01-02")+".log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}

	// Initialize loggers
	InfoLogger = log.New(logFile, "INFO: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	ErrorLogger = log.New(logFile, "ERROR: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
}
