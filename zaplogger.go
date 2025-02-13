package main

import (
	zap "go.uber.org/zap"
)

func zaplogger(){
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("This is my first log message using zap")
	logger.Error("This is my first error message using zap")

	// structured logging
	logger.Info("This is my first structured log message using zap", zap.String("Name", "John Doe"), zap.Int("Age", 23))
	logger.Error("This is my first structured error message using zap", zap.String("Name", "John Doe"), zap.Int("Age", 23))
}

func main() {
	zaplogger()
}