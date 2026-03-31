package logger

import (
	"os"

	"github.com/coderz-space/coderz.space/internal/config"
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

// Security event logging helpers

// LogSecurityEvent logs a security-related event for audit purposes
func LogSecurityEvent(eventType, userID, resource, action, result string, fields ...zap.Field) {
	if Logger != nil {
		allFields := append([]zap.Field{
			zap.String("event_type", "security"),
			zap.String("security_event", eventType),
			zap.String("user_id", userID),
			zap.String("resource", resource),
			zap.String("action", action),
			zap.String("result", result),
		}, fields...)
		Logger.Info("Security event", allFields...)
	}
}

// LogAuthenticationAttempt logs an authentication attempt
func LogAuthenticationAttempt(email, result string, fields ...zap.Field) {
	LogSecurityEvent("authentication", email, "auth", "login", result, fields...)
}

// LogAuthorizationFailure logs an authorization failure
func LogAuthorizationFailure(userID, resource, action, reason string, fields ...zap.Field) {
	allFields := append([]zap.Field{
		zap.String("reason", reason),
	}, fields...)
	LogSecurityEvent("authorization_failure", userID, resource, action, "denied", allFields...)
}

// LogDataAccess logs data access events
func LogDataAccess(userID, resourceType, resourceID, action string, fields ...zap.Field) {
	LogSecurityEvent("data_access", userID, resourceType, action, "success", append(fields, zap.String("resource_id", resourceID))...)
}

// LogCrossOrgAttempt logs cross-organization access attempts
func LogCrossOrgAttempt(userID, userOrg, targetOrg, resource string, fields ...zap.Field) {
	allFields := append([]zap.Field{
		zap.String("user_org", userOrg),
		zap.String("target_org", targetOrg),
	}, fields...)
	LogSecurityEvent("cross_org_violation", userID, resource, "access", "blocked", allFields...)
}

// LogRateLimitExceeded logs rate limit violations
func LogRateLimitExceeded(userID, endpoint string, fields ...zap.Field) {
	LogSecurityEvent("rate_limit", userID, endpoint, "request", "blocked", fields...)
}

// LogSuspiciousActivity logs suspicious activity
func LogSuspiciousActivity(userID, activityType, description string, fields ...zap.Field) {
	allFields := append([]zap.Field{
		zap.String("activity_type", activityType),
		zap.String("description", description),
	}, fields...)
	LogSecurityEvent("suspicious_activity", userID, "system", "detected", "flagged", allFields...)
}
