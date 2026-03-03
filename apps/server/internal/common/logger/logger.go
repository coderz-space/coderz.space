package logger

import (
	"os"

	"github.com/DSAwithGautam/CodeConquerers/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.Logger

func Initialize(config *config.Config) {

	consoleEncoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		TimeKey:      "ts",
		LevelKey:     "level",
		MessageKey:   "msg",
		CallerKey:    "caller",
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		EncodeLevel:  zapcore.CapitalColorLevelEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	})

	fileEncoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:       "ts",
		LevelKey:      "level",
		MessageKey:    "msg",
		CallerKey:     "caller",
		StacktraceKey: "stacktrace",
		EncodeTime:    zapcore.ISO8601TimeEncoder,
		EncodeLevel:   zapcore.LowercaseLevelEncoder,
		EncodeCaller:  zapcore.ShortCallerEncoder,
	})

	fileWriter := &lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    100,
		MaxBackups: 5,
		MaxAge:     7,
		Compress:   true,
	}

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), config.LOG_LEVEL),
		zapcore.NewCore(fileEncoder, zapcore.AddSync(fileWriter), config.FILE_LOG_LEVEL),
	)

	Logger = zap.New(
		core,
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel),
		zap.ErrorOutput(zapcore.AddSync(os.Stderr)), // handle zap internal errors
	).With(
		zap.String("service", "myapp"),
		zap.String("env", os.Getenv("APP_ENV")),
	)
}

func Sync() error {
	if Logger != nil {
		return Logger.Sync()
	}
	return nil
}

// Convenience functions for logging
func Debug(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Debug(msg, fields...)
	}
}

func Info(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Info(msg, fields...)
	}
}

func Warn(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Warn(msg, fields...)
	}
}

func Error(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Error(msg, fields...)
	}
}

func Fatal(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Fatal(msg, fields...)
	}
}

func Panic(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Panic(msg, fields...)
	}
}

// WithFields returns a logger with predefined fields
func WithFields(fields ...zap.Field) *zap.Logger {
	if Logger != nil {
		return Logger.With(fields...)
	}
	return nil
}
