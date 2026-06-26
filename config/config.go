package config

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	DatabaseURL string
	JWTSecret   string
	Environment string
	Port        string
}

var AppConfig Config

func LoadConfig() *Config {
	return &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
		Environment: os.Getenv("ENVIRONMENT"),
		Port:        os.Getenv("PORT"),
	}
}

func (c *Config) GetDB() (*gorm.DB, error) {
	dsn := c.DatabaseURL
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to database: %w", err)
	}
	return db, nil
}
