package logger

import (
	"sync"

	"github.com/upstreamcanoe/base/pkg/logger/internal"
	"go.uber.org/zap"
)

var (
	log      *zap.Logger
	cfg      *internal.Zap
	initOnce sync.Once
)

// InitLog 初始化日志
// opts 可修改日志配置, 不修改则使用默认配置. 配置详情见 logger.Option
func InitLog(opts ...Option) *zap.Logger {
	initOnce.Do(func() {
		cfg = internal.NewConfig()
		for _, opt := range opts {
			opt(cfg)
		}
		log = internal.NewZap(cfg)
		zap.ReplaceGlobals(log)
	})

	// log 实例已赋值, 无论是否接收返回的实例都可以正常使用
	return log
}

// L 日志实例
func L() *zap.Logger {
	if log == nil {
		InitLog() // 使用默认配置兜底
	}
	return log
}

// C 日志配置
func C() *internal.Zap {
	if cfg == nil {
		InitLog()
	}
	return cfg
}
