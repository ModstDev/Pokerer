package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Database DatabaseConfig
	JWT      JWTConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type JWTConfig struct {
	Secret string
	Issuer string
}

func Load() (Config, error) {
	// Unable to find .env isn't always a failure
	_ = godotenv.Load()

	cfg := Config{
		Database: DatabaseConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
		},
		JWT: JWTConfig{
			Secret: os.Getenv("JWT_SECRET"),
			Issuer: os.Getenv("JWT_ISSUER"),
		},
	}

	if cfg.Database.Host == "" {
		return Config{}, fmt.Errorf("DB_HOST is required")
	}

	if cfg.Database.Host == "" {
		return Config{}, fmt.Errorf("DB_HOST is required")
	}

	if cfg.Database.Port == "" {
		return Config{}, fmt.Errorf("DB_PORT is required")
	}

	if cfg.Database.User == "" {
		return Config{}, fmt.Errorf("DB_USER is required")
	}

	if cfg.Database.Name == "" {
		return Config{}, fmt.Errorf("DB_NAME is required")
	}

	if cfg.JWT.Secret == "" {
		return Config{}, fmt.Errorf("JWT_SECRET is required")
	}

	if cfg.JWT.Issuer == "" {
		return Config{}, fmt.Errorf("JWT_ISSUER is required")
	}

	return cfg, nil
}
