package mysql

import (
	"context"
	"time"

	"gorm.io/gorm/logger"

	logger2 "github.com/txze/wzkj-common/logger"
)

type GormLogger struct {
	SlowThreshold time.Duration
}

func NewGormLogger() *GormLogger {
	return &GormLogger{
		SlowThreshold: 200 * time.Millisecond, // 一般超过200毫秒就算慢查所以不使用配置进行更改
	}
}

var _ logger.Interface = (*GormLogger)(nil)

func (l *GormLogger) LogMode(lev logger.LogLevel) logger.Interface {
	return &GormLogger{}
}
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	logger2.FromContext(ctx).Info(msg, logger2.Any("info", data))
}
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	logger2.FromContext(ctx).Warn(msg, logger2.Any("data", data))
}
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	logger2.FromContext(ctx).Error(msg, logger2.Any("data", data))
}
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	// 获取运行时间
	elapsed := time.Since(begin)
	// 获取 SQL 语句和返回条数
	sql, rows := fc()
	// Gorm 错误时打印
	if err != nil {
		logger2.FromContext(ctx).Error(
			"SQL ERROR",
			logger2.Any("sql", sql),
			logger2.Any("rows", rows),
			logger2.Any("elapsed", elapsed),
		)
	}
	// 慢查询日志
	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		logger2.FromContext(ctx).Warn(
			"Database Slow Log",
			logger2.Any("sql", sql),
			logger2.Any("rows", rows),
			logger2.Any("elapsed", elapsed),
		)
	} else {
		logger2.FromContext(ctx).Info(
			"SQL INFO",
			logger2.Any("sql", sql),
			logger2.Any("rows", rows),
			logger2.Any("elapsed", elapsed),
		)
	}
}
