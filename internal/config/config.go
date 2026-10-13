package config

import (
	"fmt"
	"os"
)

type Config struct {
	AppEnv 		string
	HTTPPort 	string
	DBHost 		string
	DBPort 		string
	DBName 		string
	DBUser 		string
	DBPassword 	string
}

func Load() Config {
	return Config{
		AppEnv: 	getEnv("APP_ENV", "development"),
		HTTPPort: 	getEnv("HTTP_PORT", "8080"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBName:     getEnv("DB_NAME", "orchestrator"),
		DBUser:     getEnv("DB_USER", "orchestrator"),
		DBPassword: getEnv("DB_PASSWORD", "orchestrator"),
	}
}

func (c Config) DatabaseDSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		c.DBUser,
		c.DBPassword,
		c.DBHost,
		c.DBPort,
		c.DBName,
	)
}

func getEnv (key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	
	return value
}