package internal

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var config *Zap

type Zap struct {
	Level        string `mapstructure:"level" json:"level" yaml:"level"`                            // 级别
	Format       string `mapstructure:"format" json:"format" yaml:"format"`                         // 输出
	Director     string `mapstructure:"director" json:"director"  yaml:"director"`                  // 日志文件夹
	EncodeLevel  string `mapstructure:"encode-level" json:"encode-level" yaml:"encode-level"`       // 编码级
	ShowLine     bool   `mapstructure:"show-line" json:"show-line" yaml:"show-line"`                // 显示行
	LogInConsole bool   `mapstructure:"log-in-console" json:"log-in-console" yaml:"log-in-console"` // 输出控制台
	RetentionDay int    `mapstructure:"retention-day" json:"retention-day" yaml:"retention-day"`    // 日志保留天数
}

func NewConfig() *Zap {
	return &Zap{
		// debug,info,warn,error,dpanic,panic,fatal
		Level:        "info",                // 级别
		Format:       "console",             // 输出
		Director:     "logs",                // 日志文件夹
		EncodeLevel:  "CapitalLevelEncoder", // 大写无颜色日志
		ShowLine:     true,                  // 显示行号
		LogInConsole: true,                  // 输出到控制台
		RetentionDay: -1,                    // 日志保留天数, 小于等于0则永久保留
	}
}

// Levels 根据字符串转化为 zapcore.Levels
func (c *Zap) Levels() []zapcore.Level {
	levels := make([]zapcore.Level, 0, 7)
	level, err := zapcore.ParseLevel(c.Level)
	if err != nil {
		level = zapcore.DebugLevel
	}
	for ; level <= zapcore.FatalLevel; level++ {
		levels = append(levels, level)
	}
	return levels
}

// Encoder 日志编码器: 一些基础配置, 如时间格式,颜色,输出格式等
func (c *Zap) Encoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")
	encoderConfig.EncodeLevel = c.LevelEncoder()

	/*
		# ConsoleEncoder（人类友好，制表符分隔）
		2024-01-15 14:23:01.123  INFO  user/service.go:42  用户登录成功  {"uid": 123, "ip": "192.168.1.1"}

		# JSONEncoder（机器友好，结构化）
		{"ts":"2024-01-15 14:23:01.123","level":"INFO","caller":"user/service.go:42","msg":"用户登录成功","uid":123,"ip":"192.168.1.1"}
	*/
	if c.Format == "json" {
		return zapcore.NewJSONEncoder(encoderConfig)
	}
	return zapcore.NewConsoleEncoder(encoderConfig)

}

// LevelEncoder 编码级别: 大小写 + 颜色
func (c *Zap) LevelEncoder() zapcore.LevelEncoder {
	switch {
	case c.EncodeLevel == "LowercaseLevelEncoder": // 小写编码器(默认)
		return zapcore.LowercaseLevelEncoder
	case c.EncodeLevel == "LowercaseColorLevelEncoder": // 小写编码器带颜色
		return zapcore.LowercaseColorLevelEncoder
	case c.EncodeLevel == "CapitalLevelEncoder": // 大写编码器
		return zapcore.CapitalLevelEncoder
	case c.EncodeLevel == "CapitalColorLevelEncoder": // 大写编码器带颜色
		return zapcore.CapitalColorLevelEncoder
	default:
		return zapcore.CapitalLevelEncoder
	}
}
