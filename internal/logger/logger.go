package logger

import (
	"go.uber.org/zap"
)

// Initialize creates and returns a new zap logger instance
func Initialize() *zap.Logger {
	var err error
	log, err := zap.NewProduction()
	if err != nil {
		panic("failed to initialize zap logger: " + err.Error())
	}
	return log
}
