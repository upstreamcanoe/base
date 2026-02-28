package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Encoder 日志编码器: 一些基础配置, 如时间格式,颜色,输出格式等
func Encoder(encodeLevel string, isJsonEncoder bool) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")
	encoderConfig.EncodeLevel = LevelEncoder(encodeLevel)

	/*
		# ConsoleEncoder（人类友好，制表符分隔）
		2024-01-15 14:23:01.123  INFO  user/service.go:42  用户登录成功  {"uid": 123, "ip": "192.168.1.1"}

		# JSONEncoder（机器友好，结构化）
		{"ts":"2024-01-15 14:23:01.123","level":"INFO","caller":"user/service.go:42","msg":"用户登录成功","uid":123,"ip":"192.168.1.1"}
	*/
	if isJsonEncoder {
		return zapcore.NewJSONEncoder(encoderConfig)
	}
	return zapcore.NewConsoleEncoder(encoderConfig)

}

// LevelEncoder 编码级别: 大小写 + 颜色
func LevelEncoder(encodeLevel string) zapcore.LevelEncoder {
	switch {
	case encodeLevel == "LowercaseLevelEncoder": // 小写编码器
		return zapcore.LowercaseLevelEncoder
	case encodeLevel == "LowercaseColorLevelEncoder": // 小写编码器带颜色
		return zapcore.LowercaseColorLevelEncoder
	case encodeLevel == "CapitalLevelEncoder": // 大写编码器
		return zapcore.CapitalLevelEncoder
	case encodeLevel == "CapitalColorLevelEncoder":
		return zapcore.CapitalColorLevelEncoder
	default:
		return zapcore.CapitalLevelEncoder
	}
}
