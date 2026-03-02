package logger

import (
	"fmt"
	"os"

	"github.com/DSAwithGautam/CodeConquerers/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(config *config.Config) *zap.Logger {

	// ----- Encoders -----
	consoleEncoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		TimeKey:      "ts",
		LevelKey:     config.LOG_LEVEL,
		MessageKey:   "msg",
		CallerKey:    "caller",
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		EncodeLevel:  zapcore.CapitalColorLevelEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	})

	fileEncoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:       "ts",
		LevelKey:      config.LOG_LEVEL,
		MessageKey:    "msg",
		CallerKey:     "caller",
		StacktraceKey: "stacktrace",
		EncodeTime:    zapcore.ISO8601TimeEncoder,
		EncodeLevel:   zapcore.LowercaseLevelEncoder,
		EncodeCaller:  zapcore.ShortCallerEncoder,
	})

	// ----- Outputs -----
	consoleWS := zapcore.AddSync(os.Stdout)

	// -------  lumberjack ---------
	if err := os.MkdirAll("logs", 0755); err != nil {
		panic(fmt.Errorf("failed to create directories: %w", err))
	}
	file, err := os.OpenFile("logs/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	fileWS := zapcore.AddSync(file)

	// ----- Cores -----
	consoleCore := zapcore.NewCore(consoleEncoder, consoleWS, zapcore.DebugLevel)
	fileCore := zapcore.NewCore(fileEncoder, fileWS, zapcore.InfoLevel)

	// ----- Tee -----
	core := zapcore.NewTee(consoleCore, fileCore)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
	return logger
}
