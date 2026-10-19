package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Server ServerConfig
	DB     DBConfig
	JWT    JWTConfig
}

type ServerConfig struct {
	Port int
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

type JWTConfig struct {
	Secret string
}

var defConfig = map[string]string{
	"SERVER_PORT": "8080",
	"DB_PORT":     "5432",
	"DB_HOST":     "localhost",
	"DB_NAME":     "messenger",
	"DB_USER":     "postgres",
	"DB_PASSWORD": "strong_password",
}

func Load() (*Config, error) {
	port, err := strconv.Atoi(getEnv("SERVER_PORT"))
	if err != nil {
		return nil, fmt.Errorf("invalid SERVER_PORT: %w", err)
	}

	dbPort, err := strconv.Atoi(getEnv("DB_PORT"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT: %w", err)
	}

	cfg := &Config{
		Server: ServerConfig{
			Port: port,
		},
		DB: DBConfig{
			Host:     getEnv("DB_HOST"),
			Port:     dbPort,
			User:     getEnv("DB_USER"),
			Password: getEnv("DB_PASSWORD"),
			Name:     getEnv("DB_NAME"),
		},
		JWT: JWTConfig{
			Secret: os.Getenv("JWT_SECRET"),
		},
	}

	// Секрет по умолчанию не задается: подписание токенов
	// публично известным ключом сводит JWT-аутентификацию на нет.
	if cfg.JWT.Secret == "" {
		return nil, errors.New("JWT_SECRET is required")
	}

	return cfg, nil
}

func getEnv(key string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defConfig[key]
}
