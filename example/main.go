package main

import (
	"time"

	"github.com/fy138/auroralog"
)

func main() {
	// Get the singleton logger instance
	logger := auroralog.GetLogger()

	// Set log file with rotation and retention
	err := logger.SetLogFile("app.log", 7*24*time.Hour, 24*time.Hour) // Retain logs for 7 days, rotate daily
	if err != nil {
		panic(err)
	}

	// Set the log level (DEBUG, INFO, WARN, ERROR, FATAL)
	logger.SetLevel(auroralog.DEBUG)

	// Log messages at different levels
	logger.Debug("This is a debug message") // Won't be printed since log level is INFO
	logger.Info("Application started")
	logger.Warn("This is a warning")
	logger.Error("An error occurred")
}
