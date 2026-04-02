package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Environment string

const (
	Development Environment = "development"
	Test        Environment = "test"
	Production  Environment = "production"
)

func (e Environment) isValid() bool {
	return e == Development || e == Test || e == Production
}

func getEnvVariable(key string) string {
	envVar := os.Getenv(key)
	if envVar == "" {
		panic(fmt.Sprintf("environment variable %s is not set", key))
	}
	return envVar
}

const envFilePath = ".env"

type Config struct {
	AppName             string
	Version             string
	Port                string
	JWTSecret           string
	JWTExpires          string
	FrontendOrigin      string
	DBURL               string
	Environment         Environment
	RefreshTokenExpires time.Duration
	MaxDBConnLifetime   time.Duration
	MaxDBConnIdleTime   time.Duration
	MaxDBConns          int
	MinDBConns          int
	LogLevel            zapcore.Level
	FileLogLevel        zapcore.Level
	SMTPHost            string
	SMTPPort            int
	SMTPUser            string
	SMTPPass            string
	SMTPFrom            string
}

func parseLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "dpanic":
		return zap.DPanicLevel
	case "panic":
		return zap.PanicLevel
	case "fatal":
		return zap.FatalLevel
	default:
		return zap.InfoLevel
	}
}

func LoadConfig() *Config {
	if err := godotenv.Load(envFilePath); err != nil && !errors.Is(err, os.ErrNotExist) {
		panic(fmt.Errorf("failed to load environment variables: %v", err))
	}

	// DB configs
	maxDBConns, err := strconv.Atoi(getEnvVariable("MAX_DB_CONNS"))
	if err != nil {
		panic(fmt.Errorf("failed to parse MAX_DB_CONNS: %v", err))
	}
	minDBConns, err := strconv.Atoi(getEnvVariable("MIN_DB_CONNS"))
	if err != nil {
		panic(fmt.Errorf("failed to parse MIN_DB_CONNS: %v", err))
	}
	maxDBConnLifetime, err := time.ParseDuration(getEnvVariable("MAX_DB_CONN_LIFETIME"))
	if err != nil {
		panic(fmt.Errorf("failed to parse MAX_DB_CONN_LIFETIME: %v", err))
	}
	maxDBConnIdleTime, err := time.ParseDuration(getEnvVariable("MAX_DB_CONN_IDLE_TIME"))
	if err != nil {
		panic(fmt.Errorf("failed to parse MAX_DB_CONN_IDLE_TIME: %v", err))
	}
	refreshTokenExpires, err := time.ParseDuration(getEnvVariable("REFRESH_TOKEN_EXPIRES"))
	if err != nil {
		panic(fmt.Errorf("failed to parse REFRESH_TOKEN_EXPIRES: %v", err))
	}
	smtpPort, _ := strconv.Atoi(os.Getenv("SMTP_PORT")) // Optional, defaults to 0 if not set or invalid
	if smtpPort == 0 {
		smtpPort = 587 // Default SMTP port
	}

	config := &Config{

		AppName:             getEnvVariable("APP_NAME"),
		Version:             getEnvVariable("VERSION"),
		Environment:         Environment(getEnvVariable("ENVIRONMENT")),
		Port:                getEnvVariable("PORT"),
		JWTSecret:           getEnvVariable("JWT_SECRET"),
		JWTExpires:          getEnvVariable("JWT_EXPIRES"),
		LogLevel:            parseLevel(getEnvVariable("LOG_LEVEL")),
		FileLogLevel:        parseLevel(getEnvVariable("FILE_LOG_LEVEL")),
		FrontendOrigin:      getEnvVariable("FRONTEND_ORIGIN"),
		DBURL:               getEnvVariable("DB_URL"),
		MaxDBConns:          maxDBConns,
		MinDBConns:          minDBConns,
		MaxDBConnLifetime:   maxDBConnLifetime,
		MaxDBConnIdleTime:   maxDBConnIdleTime,
		RefreshTokenExpires: refreshTokenExpires,
		SMTPHost:            os.Getenv("SMTP_HOST"),
		SMTPPort:            smtpPort,
		SMTPUser:            os.Getenv("SMTP_USER"),
		SMTPPass:            os.Getenv("SMTP_PASS"),
		SMTPFrom:            os.Getenv("SMTP_FROM"),
	}

	if !config.Environment.isValid() {
		panic(fmt.Errorf("invalid environment: %s", config.Environment))
	}

	return config
}
