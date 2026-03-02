package logger

import (
	"github.com/upstreamcanoe/base/pkg/logger/internal"
)

// Option 日志配置项
type Option func(*internal.Zap)

// WithLevel 设置日志级别
// 默认 info, 可选 debug, info, warn, error, dpanic, panic, fatal
func WithLevel(level string) Option {
	return func(z *internal.Zap) {
		z.Level = level
	}
}

// WithPrefix 设置日志前缀
// 默认 ""
//func WithPrefix(prefix string) Option {
//	return func(z *internal.Zap) {
//		z.Prefix = prefix
//	}
//}

// WithFormat 输出格式
// 默认 console
func WithFormat(format string) Option {
	return func(z *internal.Zap) {
		z.Format = format
	}
}

// WithDirector 日志目录
// 默认 "logs"
func WithDirector(director string) Option {
	return func(z *internal.Zap) {
		z.Director = director
	}
}

// WithEncodeLevel 日志编码级别
// 默认 LowercaseColorLevelEncoder
// 可选值: LowercaseLevelEncoder, LowercaseColorLevelEncoder, CapitalLevelEncoder, CapitalColorLevelEncoder
func WithEncodeLevel(encodeLevel string) Option {
	return func(z *internal.Zap) {
		z.EncodeLevel = encodeLevel
	}
}

// WithStacktraceKey 栈名
// 默认 ""
//func WithStacktraceKey(stacktraceKey string) Option {
//	return func(z *internal.Zap) {
//		z.StacktraceKey = stacktraceKey
//	}
//}

// WithShowLine 显示行号
// 默认 true
func WithShowLine(showLine bool) Option {
	return func(z *internal.Zap) {
		z.ShowLine = showLine
	}
}

// WithLogInConsole 是否输出到控制台
// 默认 true
func WithLogInConsole(logInConsole bool) Option {
	return func(z *internal.Zap) {
		z.LogInConsole = logInConsole
	}
}

// WithRetentionDay 日志保留天数
// 默认 -1, 小于等于0为永久保留
func WithRetentionDay(retentionDay int) Option {
	return func(z *internal.Zap) {
		z.RetentionDay = retentionDay
	}
}
