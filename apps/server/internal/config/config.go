package config

import (
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
	AppName        string
	Version        string
	Environment    Environment
	Port           string
	JWT_SECRET     string
	JWT_EXPIRES    string
	LOG_LEVEL      zapcore.Level
	FILE_LOG_LEVEL zapcore.Level
	FrontendOrigin string
	// DB config
	DB_URL            string
	MaxDBConns        int
	MinDBConns        int
	MaxDBConnLifetime time.Duration
	MaxDBConnIdleTime time.Duration
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
	if err := godotenv.Load(envFilePath); err != nil {
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

	config := &Config{

		AppName:           getEnvVariable("APP_NAME"),
		Version:           getEnvVariable("VERSION"),
		Environment:       Environment(getEnvVariable("ENVIRONMENT")),
		Port:              getEnvVariable("PORT"),
		JWT_SECRET:        getEnvVariable("JWT_SECRET"),
		JWT_EXPIRES:       getEnvVariable("JWT_EXPIRES"),
		LOG_LEVEL:         parseLevel(getEnvVariable("LOG_LEVEL")),
		FILE_LOG_LEVEL:    parseLevel(getEnvVariable("FILE_LOG_LEVEL")),
		FrontendOrigin:    getEnvVariable("FRONTEND_ORIGIN"),
		DB_URL:            getEnvVariable("DB_URL"),
		MaxDBConns:        maxDBConns,
		MinDBConns:        minDBConns,
		MaxDBConnLifetime: maxDBConnLifetime,
		MaxDBConnIdleTime: maxDBConnIdleTime,
	}

	if !config.Environment.isValid() {
		panic(fmt.Errorf("invalid environment: %s", config.Environment))
	}

	return config
}
