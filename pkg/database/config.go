package database

import (
	"github.com/serediukit/civix-backend/pkg/env"
)

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

func GetDBConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Host:     env.GetEnv("DB_HOST", "localhost"),
		Port:     env.GetEnv("DB_PORT", "5432"),
		User:     env.GetEnv("DB_USER", "db"),
		Password: env.GetEnv("DB_PASSWORD", "db"),
		Name:     env.GetEnv("DB_NAME", "civix_db"),
		SSLMode:  env.GetEnv("DB_SSLMODE", "disable"),
	}
}
