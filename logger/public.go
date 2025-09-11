package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var hookCore []zapcore.Core

type L struct {
	log *zap.Logger
}

func (l *L) Debug(msg string, field ...zap.Field) {
	l.log.Debug(msg, field...)
}

func (l *L) Warn(msg string, field ...zap.Field) {
	l.log.Warn(msg, field...)
}

func (l *L) Info(msg string, field ...zap.Field) {
	l.log.Info(msg, field...)
}

func (l *L) Error(msg string, field ...zap.Field) {
	l.log.Error(msg, field...)
}

func (l *L) Fatal(msg string, field ...zap.Field) {
	l.log.Fatal(msg, field...)
}

//func (l *L) Errorf(format string, f ...interface{}) {
//	defaultLogger.log.Error(fmt.Sprintf(format, f...))
//}
//
//func (l *L) Warnf(format string, f ...interface{}) {
//	defaultLogger.log.Warn(fmt.Sprintf(format, f...))
//}
//
//func (l *L) Infof(format string, f ...interface{}) {
//	defaultLogger.log.Info(fmt.Sprintf(format, f...))
//}
//
//func (l *L) Debugf(format string, f ...interface{}) {
//	defaultLogger.log.Debug(fmt.Sprintf(format, f...))
//}

//func DebugAsJson(value interface{}) {
//	defaultLogger.Debug("debugAsJson", zap.Any("object", value))
//}

func Err(err error) zap.Field {
	return zap.Error(err)
}

func String(key string, val string) zap.Field {
	return zap.String(key, val)
}

func Any(key string, value interface{}) zap.Field {
	return zap.Any(key, value)
}

func Binary(key string, val []byte) zap.Field {
	return zap.Binary(key, val)
}

func Bool(key string, val bool) zap.Field {
	return zap.Bool(key, val)
}

func ByteString(key string, val []byte) zap.Field {
	return zap.ByteString(key, val)
}

func Float64(key string, val float64) zap.Field {
	return zap.Float64(key, val)
}

func Float32(key string, val float32) zap.Field {
	return zap.Float32(key, val)
}

func NewEncoderConfig() zapcore.EncoderConfig {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.CallerKey = "line"
	encoderConfig.TimeKey = "time"
	encoderConfig.StacktraceKey = ""
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	return encoderConfig
}

func newHook(core zapcore.Core) {
	hookCore = []zapcore.Core{
		core,
	}
}

// 针对这个日志可以配置输出位置
func consoleHook(path string) zapcore.Core {
	encoderConfig := NewEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	level := zap.LevelEnablerFunc(func(lv zapcore.Level) bool {
		return lv <= zap.FatalLevel
	})
	// 设置日志输出
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel && lvl >= zapcore.DebugLevel
	})

	consoleWriter := zapcore.AddSync(os.Stdout)
	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)

	consoleCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		consoleDebugging,
		lowPriority,
	)

	consoleErrorCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		consoleErrors,
		highPriority,
	)

	fileCore := zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), consoleWriter, level)

	core := zapcore.NewTee(consoleCore, consoleErrorCore, fileCore)

	return core

}

func New(path string) *L {
	newHook(consoleHook(path))
	return &L{log: zap.New(zapcore.NewTee(hookCore...), zap.AddCaller(), zap.AddCallerSkip(2))}

}

func Int(key string, val int) zap.Field {
	return Int64(key, int64(val))
}

func Int64(key string, val int64) zap.Field {
	return zap.Int64(key, val)
}

func Int8(key string, val int8) zap.Field {
	return zap.Int8(key, val)
}

func Uint(key string, val uint) zap.Field {
	return Uint64(key, uint64(val))
}

func Uint64(key string, val uint64) zap.Field {
	return zap.Uint64(key, val)
}

func Uint8(key string, val uint8) zap.Field {
	return zap.Uint8(key, val)
}

func SetLevel(level string) zap.AtomicLevel {
	switch level {
	case "debug":
		return zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		return zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		return zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		return zap.NewAtomicLevelAt(zap.ErrorLevel)
	case "fatal":
		return zap.NewAtomicLevelAt(zap.FatalLevel)
	case "panic":
		return zap.NewAtomicLevelAt(zap.PanicLevel)
	default:
		return zap.NewAtomicLevelAt(zap.DebugLevel)
	}
}
