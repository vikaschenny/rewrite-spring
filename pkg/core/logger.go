package core

import (
	"log"
	"os"
)

// ConsoleLogger implements the Logger interface for console output
type ConsoleLogger struct {
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
	debugLogger *log.Logger
	debugMode   bool
}

// NewConsoleLogger creates a new console logger
func NewConsoleLogger(debugMode bool) *ConsoleLogger {
	return &ConsoleLogger{
		infoLogger:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime),
		warnLogger:  log.New(os.Stdout, "WARN: ", log.Ldate|log.Ltime),
		errorLogger: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime),
		debugLogger: log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime),
		debugMode:   debugMode,
	}
}

// Info logs an info message
func (l *ConsoleLogger) Info(msg string, args ...interface{}) {
	l.infoLogger.Printf(msg, args...)
}

// Warn logs a warning message
func (l *ConsoleLogger) Warn(msg string, args ...interface{}) {
	l.warnLogger.Printf(msg, args...)
}

// Error logs an error message
func (l *ConsoleLogger) Error(msg string, args ...interface{}) {
	l.errorLogger.Printf(msg, args...)
}

// Debug logs a debug message (only if debug mode is enabled)
func (l *ConsoleLogger) Debug(msg string, args ...interface{}) {
	if l.debugMode {
		l.debugLogger.Printf(msg, args...)
	}
}

// NullLogger implements the Logger interface but does nothing (for testing)
type NullLogger struct{}

// NewNullLogger creates a new null logger
func NewNullLogger() *NullLogger {
	return &NullLogger{}
}

// Info does nothing
func (l *NullLogger) Info(msg string, args ...interface{}) {}

// Warn does nothing
func (l *NullLogger) Warn(msg string, args ...interface{}) {}

// Error does nothing
func (l *NullLogger) Error(msg string, args ...interface{}) {}

// Debug does nothing
func (l *NullLogger) Debug(msg string, args ...interface{}) {}
