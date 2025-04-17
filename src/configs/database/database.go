package database

import (
	"fmt"
	"gin/src/entities/users"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found or could not be loaded")
	}

	host := os.Getenv("HOST")
	portStr := os.Getenv("PORT")
	user := os.Getenv("USER")
	dbname := os.Getenv("DB_NAME")
	pass := os.Getenv("PASSWORD")

	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Invalid port number: %v", err)
	}

	// PostgreSQL DSN
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s sslmode=disable TimeZone=UTC",
		host, port, user, dbname, pass,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database using GORM: %v", err)
	}

	// Test connection
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Database is not reachable: %v", err)
	}

	// Configure connection pool
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Set global DB instance
	DB = db
	fmt.Println("Successfully connected to database using GORM!")

	// Call CheckTables to check and migrate models
	CheckTables(&users.User{})
}

// CheckTables checks if the tables exist for all provided models and migrates them if not.
func CheckTables(models ...interface{}) {
	for _, model := range models {
		if !DB.Migrator().HasTable(model) {
			if err := DB.AutoMigrate(model); err != nil {
				log.Fatalf("Auto migration failed for %v: %v", model, err)
			}
			fmt.Printf("%v table migration completed successfully!\n", model)
		} else {
			fmt.Printf("%v table already exists, skipping migration.\n", model)
		}
	}
}
