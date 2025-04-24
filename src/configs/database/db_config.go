package database

import (
	"fmt"
	"os"
	"strconv"
)

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
	Timezone string
}

func LoadDBConfig() DBConfig {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		fmt.Println("Invalid port number: %v", err)
	}

	tz := os.Getenv("TIMEZONE")
	if tz == "" {
		tz = "UTC"
	}

	return DBConfig{
		Host:     os.Getenv("HOST"),
		Port:     port,
		User:     os.Getenv("USER"),
		Password: os.Getenv("PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("SSL_MODE"),
		Timezone: tz,
	}
}
func (cfg DBConfig) ToDSN() string {
	if cfg.Password == "" {
		return fmt.Sprintf(
			"host=%s port=%d user=%s dbname=%s sslmode=%s TimeZone=%s",
			cfg.Host, cfg.Port, cfg.User, cfg.DBName, cfg.SSLMode, cfg.Timezone,
		)
	}

	return fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s sslmode=%s TimeZone=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.DBName, cfg.Password, cfg.SSLMode, cfg.Timezone,
	)
}
