package logger

import (
	"context"

	"github.com/gin-gonic/gin"
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

// WithTrace 创建带追踪ID的日志记录
func WithTrace(c *gin.Context) *zap.Logger {
	traceID := c.GetString("traceID")
	if traceID == "" {
		traceID = "unknown"
	}

	return defaultLogger.log.With(
		zap.String("traceID", traceID),
		zap.String("path", c.Request.URL.Path),
		zap.String("method", c.Request.Method),
		zap.String("ip", c.ClientIP()),
	)
}

// FromContext 从context获取带追踪ID的日志记录
func FromContext(ctx context.Context) *zap.Logger {
	traceID := ctx.Value("traceID")
	if traceID == nil {
		traceID = "unknown"
	}

	return defaultLogger.log.With(
		zap.String("traceID", traceID.(string)),
	)
}
