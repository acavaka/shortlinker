package logger

import (
	"go.uber.org/zap"
)

var log *zap.Logger

// Initialize creates and returns a new zap logger instance
func Initialize() *zap.Logger {
	var err error
	log, err = zap.NewProduction()
	if err != nil {
		panic("failed to initialize zap logger: " + err.Error())
	}
	return log
}

// Error logs an error message with additional fields
func Error(msg string, err interface{}) {
	if log == nil {
		log = Initialize()
	}
	log.Error(msg, zap.Any("error", err))
}

// Info logs an info message with additional fields
func Info(msg string, fields ...zap.Field) {
	if log == nil {
		log = Initialize()
	}
	log.Info(msg, fields...)
}
