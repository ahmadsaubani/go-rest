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

// LoadDBConfig reads the environment variables and returns a DBConfig instance.
// The port number is expected to be an integer, and the timezone is expected to
// be a valid string. If the port number is invalid, a warning message is printed
// to the console. If the timezone is empty, it defaults to "UTC".
func LoadDBConfig() DBConfig {
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		fmt.Println("invalid port number: %w", err)
	}

	tz := os.Getenv("DB_TIMEZONE")
	if tz == "" {
		tz = "UTC"
	}

	config := DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     port,
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSL_MODE"),
		Timezone: tz,
	}

	// for debugging purposes, you can uncomment the following lines to print the configuration
	// fmt.Println("=== Database Configuration ===")
	// fmt.Printf("Host: %s\n", config.Host)
	// fmt.Printf("Port: %d\n", config.Port)
	// fmt.Printf("User: %s\n", config.User)
	// fmt.Printf("Password: %s\n", config.Password)
	// fmt.Printf("DB Name: %s\n", config.DBName)
	// fmt.Printf("SSL Mode: %s\n", config.SSLMode)
	// fmt.Printf("Timezone: %s\n", config.Timezone)
	// fmt.Println("===============================")

	return config
}

// ToDSN converts the DBConfig instance to a PostgreSQL connection string (DSN).
// The connection string is a URL-encoded string that includes the host, port, user,
// password, database name, SSL mode, and timezone.
// If the password is empty, it is omitted from the connection string.
// This method is intended to be used for creating a connection to a PostgreSQL
// database using the github.com/lib/pq package.
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
