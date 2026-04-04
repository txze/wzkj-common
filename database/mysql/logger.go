package mysql

import (
	"context"
	"strings"
	"time"

	"go.uber.org/zap"
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
	logger2.FromContext(ctx).Info(msg, logger2.Any("data", data))
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

	// 提取 SQL 操作类型
	opType := extractSQLType(sql)

	// 构建日志字段
	fields := []zap.Field{
		logger2.String("sql", sql),
		logger2.Int64("rows", rows),
		logger2.String("elapsed", elapsed.String()),
		logger2.String("op_type", opType),
	}

	// Gorm 错误时打印
	if err != nil {
		fields = append(fields, logger2.Err(err))
		logger2.FromContext(ctx).Error("SQL ERROR", fields...)
		return
	}

	// 慢查询日志
	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		logger2.FromContext(ctx).Warn("DATABASE SLOW QUERY", fields...)
	} else {
		logger2.FromContext(ctx).Info("SQL EXECUTION", fields...)
	}
}

// extractSQLType 提取 SQL 操作类型
func extractSQLType(sql string) string {
	// 简单的 SQL 类型提取，实际项目中可能需要更复杂的解析
	sqlUpper := strings.ToUpper(sql)
	if strings.HasPrefix(sqlUpper, "SELECT") {
		return "SELECT"
	} else if strings.HasPrefix(sqlUpper, "INSERT") {
		return "INSERT"
	} else if strings.HasPrefix(sqlUpper, "UPDATE") {
		return "UPDATE"
	} else if strings.HasPrefix(sqlUpper, "DELETE") {
		return "DELETE"
	} else if strings.HasPrefix(sqlUpper, "CREATE") {
		return "CREATE"
	} else if strings.HasPrefix(sqlUpper, "DROP") {
		return "DROP"
	} else if strings.HasPrefix(sqlUpper, "ALTER") {
		return "ALTER"
	}
	return "OTHER"
}
