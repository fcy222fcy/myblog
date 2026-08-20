package logger

import (
	"io"
	"os"
	"path/filepath"

	"blog/pkg/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	log          = zap.NewNop()
	sugar        = log.Sugar()
	outputCloser io.Closer
)

// Init 初始化日志系统
func Init(cfg config.LogConfig) error {
	_ = closeOutputs()

	if err := ensureLogDir(cfg.Filename); err != nil {
		return err
	}

	level := parseLevel(cfg.Level)

	fileWriter := &lumberjack.Logger{
		Filename:   cfg.Filename,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
	}
	outputCloser = fileWriter

	core := zapcore.NewTee(
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(newConsoleEncoderConfig()),
			zapcore.AddSync(os.Stdout),
			level,
		),
		zapcore.NewCore(
			zapcore.NewJSONEncoder(newJSONEncoderConfig()),
			zapcore.AddSync(fileWriter),
			level,
		),
	)

	log = zap.New(
		core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	sugar = log.Sugar()
	return nil
}

func parseLevel(raw string) zapcore.Level {
	switch raw {
	case "debug":
		return zapcore.DebugLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

func newBaseEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func newConsoleEncoderConfig() zapcore.EncoderConfig {
	cfg := newBaseEncoderConfig()
	cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	cfg.ConsoleSeparator = "  "
	return cfg
}

func newJSONEncoderConfig() zapcore.EncoderConfig {
	cfg := newBaseEncoderConfig()
	cfg.EncodeLevel = zapcore.LowercaseLevelEncoder
	return cfg
}

func ensureLogDir(filename string) error {
	dir := filepath.Dir(filename)
	if dir == "" || dir == "." {
		return nil
	}
	return os.MkdirAll(dir, 0o755)
}

func closeOutputs() error {
	if outputCloser == nil {
		return nil
	}
	err := outputCloser.Close()
	outputCloser = nil
	return err
}

func Debug(msg string, fields ...zap.Field) { log.Debug(msg, fields...) }
func Info(msg string, fields ...zap.Field)  { log.Info(msg, fields...) }
func Warn(msg string, fields ...zap.Field)  { log.Warn(msg, fields...) }
func Error(msg string, fields ...zap.Field) { log.Error(msg, fields...) }
func Fatal(msg string, fields ...zap.Field) { log.Panic(msg, fields...) }

func Debugf(format string, args ...interface{}) { sugar.Debugf(format, args...) }
func Infof(format string, args ...interface{})  { sugar.Infof(format, args...) }
func Warnf(format string, args ...interface{})  { sugar.Warnf(format, args...) }
func Errorf(format string, args ...interface{}) { sugar.Errorf(format, args...) }
func Fatalf(format string, args ...interface{}) { sugar.Fatalf(format, args...) }

func With(fields ...zap.Field) *zap.Logger { return log.With(fields...) }
func Sync() error                          { return log.Sync() }
func GetLogger() *zap.Logger               { return log }
func GetSugaredLogger() *zap.SugaredLogger { return sugar }
