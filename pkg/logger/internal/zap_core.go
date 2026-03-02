package internal

import (
	"fmt"
	"os"
	"time"

	"github.com/upstreamcanoe/base/pkg/utils"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewZap(c *Zap) (logger *zap.Logger) {
	config = c

	// 判断是否有 Director 文件夹
	if ok := utils.PathExists(config.Director); !ok {
		fmt.Printf("create %v directory\n", config.Director)
		_ = os.Mkdir(config.Director, os.ModePerm)
	}
	levels := config.Levels()
	length := len(levels)
	cores := make([]zapcore.Core, 0, length)
	for i := 0; i < length; i++ {
		core := NewZapCore(levels[i])
		cores = append(cores, core)
	}
	// 构建基础 logger（错误级别的入库逻辑已在自定义 ZapCore 中处理）
	logger = zap.New(zapcore.NewTee(cores...))
	// 启用 Error 及以上级别的堆栈捕捉，确保 entry.Stack 可用
	opts := []zap.Option{zap.AddStacktrace(zapcore.ErrorLevel)}
	if config.ShowLine {
		opts = append(opts, zap.AddCaller())
	}
	logger = logger.WithOptions(opts...)
	return logger
}

type ZapCore struct {
	level zapcore.Level
	zapcore.Core
}

func NewZapCore(level zapcore.Level) *ZapCore {
	entity := &ZapCore{level: level}
	syncer := entity.WriteSyncer()
	levelEnabler := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l == level
	})
	entity.Core = zapcore.NewCore(config.Encoder(), syncer, levelEnabler)
	return entity
}

func (z *ZapCore) WriteSyncer(formats ...string) zapcore.WriteSyncer {
	cutter := NewCutter(
		config.Director,
		z.level.String(),
		config.RetentionDay,
		CutterWithLayout(time.DateOnly),
		CutterWithFormats(formats...),
	)
	if config.LogInConsole {
		multiSyncer := zapcore.NewMultiWriteSyncer(os.Stdout, cutter)
		return zapcore.AddSync(multiSyncer)
	}

	return zapcore.AddSync(cutter)
}

func (z *ZapCore) Enabled(level zapcore.Level) bool {
	return z.level == level
}

func (z *ZapCore) With(fields []zapcore.Field) zapcore.Core {
	return z.Core.With(fields)
}

func (z *ZapCore) Check(entry zapcore.Entry, check *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if z.Enabled(entry.Level) {
		return check.AddCore(entry, z)
	}
	return check
}

func (z *ZapCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	return z.Core.Write(entry, fields)
}

func (z *ZapCore) Sync() error {
	return z.Core.Sync()
}
