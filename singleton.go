// Question: Implement a logging system that ensures only one instance of the logger is used throughout the application.
// Scenario: Create a Logger class that provides a global point of access to the logging instance.

package main

import (
	"fmt"
	"sync"
	"time"
)

type LogLevel int

const (
	INFO LogLevel = iota
	WARNING
	ERROR
)

// Logger struct
type Logger struct {
	message string
	mu      sync.Mutex
}

var instance *Logger
var once sync.Once

// GetLogger function - thread-safe singleton
func GetLogger() *Logger {
	once.Do(func() {
		instance = &Logger{}
	})
	return instance
}

// LogMessage function with timestamp and log level
func (l *Logger) LogMessage(level LogLevel) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.message == "" {
		return fmt.Errorf("empty message")
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	levelStr := "INFO"
	switch level {
	case WARNING:
		levelStr = "WARNING"
	case ERROR:
		levelStr = "ERROR"
	}

	fmt.Printf("[%s] %s: %s\n", timestamp, levelStr, l.message)
	return nil
}

// SetMessage function with thread safety
func (l *Logger) SetMessage(message string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.message = message
}

func main() {
	logger := GetLogger()
	logger.SetMessage("Hello, World!")
	logger.LogMessage(INFO)

	// Test singleton
	logger2 := GetLogger()
	logger2.SetMessage("Testing singleton")
	logger2.LogMessage(WARNING)

	// Verify both loggers are the same instance
	fmt.Printf("Same instance: %v\n", logger == logger2)
}
