package logger

import (
	"go.uber.org/zap"
)

var defaultLogger *L

func Debug(msg string, field ...zap.Field) {
	defaultLogger.Debug(msg, field...)
}

func Warn(msg string, field ...zap.Field) {
	defaultLogger.Warn(msg, field...)
}

func Info(msg string, field ...zap.Field) {
	defaultLogger.Info(msg, field...)
}

func Error(msg string, field ...zap.Field) {
	defaultLogger.Error(msg, field...)
}

func Fatal(msg string, field ...zap.Field) {
	defaultLogger.Fatal(msg, field...)
}
