package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
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
	Environment Environment
	Port        string
	DB_URL      string
	JWT_SECRET  string
	LOG_LEVEL   string
}

func LoadConfig() *Config {
	if err := godotenv.Load(envFilePath); err != nil {
		panic(fmt.Errorf("failed to load environment variables: %v", err))
	}

	config := &Config{
		Environment: Environment(getEnvVariable("ENVIRONMENT")),
		Port:        getEnvVariable("PORT"),
		DB_URL:      getEnvVariable("DB_URL"),
		JWT_SECRET:  getEnvVariable("JWT_SECRET"),
		LOG_LEVEL:   getEnvVariable("LOG_LEVEL"),
	}

	if !config.Environment.isValid() {
		panic(fmt.Errorf("invalid environment: %s", config.Environment))
	}

	return config
}
