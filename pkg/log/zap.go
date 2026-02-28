package log

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/natefinch/lumberjack"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm/logger"
)

var Log *zap.Logger

type ZapStarter struct {
}

func (s *ZapStarter) Init() error {
	Init()
	return nil
}

func Init() {
	// 编码配置: 颜色/路径/时间/输出格式
	encoder := Encoder(viper.GetString("log.encode-level"), viper.GetString("env") != "dev")

	// 级别配置: info+warn走log.normal; error+panic+fatal走log.error
	normalPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool { // info 级别
		return lev < zap.ErrorLevel && lev >= zap.InfoLevel
	})
	errorPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool { // error 级别
		return lev >= zap.ErrorLevel
	})

	Log = zap.New(zapcore.NewTee(
		[]zapcore.Core{
			zapcore.NewCore(
				encoder,
				zapcore.NewMultiWriteSyncer(newRollingFile("log.normal"), zapcore.AddSync(os.Stdout)),
				normalPriority,
			),
			zapcore.NewCore(
				encoder,
				zapcore.NewMultiWriteSyncer(newRollingFile("log.error"), zapcore.AddSync(os.Stdout)),
				errorPriority,
			),
		}...,
	), zap.AddCaller(), zap.AddCallerSkip(1))
}

func newRollingFile(prefix string) zapcore.WriteSyncer {
	return zapcore.AddSync(&lumberjack.Logger{
		Filename:   viper.GetString(prefix + ".file"),
		MaxSize:    viper.GetInt(prefix + ".size"),
		MaxBackups: viper.GetInt(prefix + ".backups"),
		MaxAge:     viper.GetInt(prefix + ".age"),
		Compress:   viper.GetBool(prefix + ".compress"),
		LocalTime:  true,
	})
}

func traceID(ctx context.Context, fields ...zap.Field) []zap.Field {
	sp := trace.SpanContextFromContext(ctx)
	traceID := zap.String("trace_id", sp.TraceID().String())
	fields = append(fields, traceID)
	return fields
}

func InfoWithCtx(ctx context.Context, msg string, fields ...zap.Field) {
	Info(msg, traceID(ctx, fields...)...)
}

func Info(msg string, fields ...zap.Field) {
	if strings.Contains(msg, "SHOW STATUS") {
		return
	}
	if strings.Contains(msg, "GET _health") {
		return
	}
	if strings.Contains(msg, "GET url:  /health") {
		return
	}
	Log.Info(msg, fields...)
}

func WarnWithCtx(ctx context.Context, msg string, fields ...zap.Field) {
	Warn(msg, traceID(ctx, fields...)...)
}

func Warn(msg string, fields ...zap.Field) {
	Log.Warn(msg, fields...)
}

func ErrorWithCtx(ctx context.Context, msg string, fields ...zap.Field) {
	Log.Error(msg, traceID(ctx, fields...)...)
}

func Error(msg string, fields ...zap.Field) {
	Log.Error(msg, fields...)
}

func InfoFWithCtx(ctx context.Context, msg string, args ...any) {
	Log.Info(fmt.Sprintf(msg, args...), traceID(ctx)...)
}

func InfoF(msg string, args ...any) {
	msg = fmt.Sprintf(msg, args...)
	Info(msg)
}

func WarnFWithCtx(ctx context.Context, msg string, args ...any) {
	Log.Warn(fmt.Sprintf(msg, args...), traceID(ctx)...)
}

func WarnF(msg string, args ...any) {
	Log.Warn(fmt.Sprintf(msg, args...))
}

func ErrorFWithCtx(ctx context.Context, msg string, args ...any) {
	Log.Warn(fmt.Sprintf(msg, args...), traceID(ctx)...)
}

func ErrorF(msg string, args ...any) {
	Log.Error(fmt.Sprintf(msg, args...))
}

type GLog struct {
}

func (g GLog) LogMode(_ logger.LogLevel) logger.Interface {
	return &g
}

func (g GLog) Info(_ context.Context, msg string, args ...any) {
	Info(fmt.Sprintf(msg, args...))
}

func (g GLog) Warn(_ context.Context, msg string, args ...any) {
	Warn(fmt.Sprintf(msg, args...))
}

func (g GLog) Error(_ context.Context, msg string, args ...any) {
	Error(fmt.Sprintf(msg, args...))
}

func (g GLog) Trace(_ context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, ra := fc()
	if sql == "SHOW STATUS" {
		return
	}

	sqlSpaceReg := regexp.MustCompile(`\s+`)
	formatSql := strings.TrimSpace(sqlSpaceReg.ReplaceAllString(sql, " "))
	Info("gorm callback trace", zap.String("sql", formatSql), zap.Int64("rows", ra), zap.Error(err))
}

func (g GLog) Printf(msg string, args ...any) {
	Info(fmt.Sprintf(msg, args...))
}

// CLog cron logger
type CLog struct {
}

func (c CLog) Info(msg string, keysAndValues ...any) {
	InfoF(msg, keysAndValues...)
}

// Error logs an error condition.
func (c CLog) Error(err error, msg string, keysAndValues ...any) {
	ErrorF(msg, keysAndValues...)
}

func (c CLog) Printf(msg string, args ...any) {
	InfoF(msg, args...)
}
