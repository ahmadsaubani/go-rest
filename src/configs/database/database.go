package database

import (
	"database/sql"
	"fmt"
	"gin/src/entities/auth"
	"gin/src/entities/users"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBConnection struct {
	Gorm *gorm.DB
	SQL  *sql.DB
}

var GormDB *gorm.DB
var SQLDB *sql.DB

func ConnectDatabase() *DBConnection {
	fmt.Println("===== Connect To Database =====")

	// Load .env
	if err := godotenv.Load(); err != nil {
		fmt.Println("⚠️ Warning: .env file not found:%w", err)
	}

	useGorm := os.Getenv("USE_GORM") == "true"

	if useGorm {
		GormDB := ConnectDatabaseUsingGorm()
		ResetDBUsingGorm(GormDB)
		return &DBConnection{Gorm: GormDB}
	} else {
		SQLDB := connectWithSQL()
		return &DBConnection{SQL: SQLDB}
	}
}

func ConnectDatabaseUsingGorm() *gorm.DB {
	fmt.Println("=====Connect To Database=====")
	// Load environment variables
	cfg := LoadDBConfig()
	dsn := cfg.ToDSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("❌ Failed to connect to database using GORM: %v", err)
	}

	// Test connection
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println("❌ Failed to get database instance: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		fmt.Println("❌ Database is not reachable: %v", err)
	}

	// Configure connection pool
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Set global DB instance
	GormDB = db
	fmt.Println("✅ Successfully connected to database using GORM!")

	// Call CheckTables to check and migrate models
	ResetDBUsingGorm(GormDB)

	return GormDB
}

// ResetDB drops and recreates all tables
func ResetDBUsingGorm(db *gorm.DB) {
	fmt.Println("=====Process Migrate all tables=====")
	fmt.Println("⚠️ Dropping all tables....")
	err := db.Migrator().DropTable(
		&users.User{},
		&auth.AccessToken{},
		&auth.RefreshToken{},
	)
	if err != nil {
		fmt.Println("❌ Failed to drop tables: %v", err)
	}

	fmt.Println("✅ Dropped all tables")

	fmt.Println("🔧 Migrating tables....")
	err = db.AutoMigrate(
		&users.User{},
		&auth.AccessToken{},
		&auth.RefreshToken{},
	)
	if err != nil {
		fmt.Println("❌ Failed to migrate tables: %v", err)
	}

	fmt.Println("✅ Database migrated successfully")
}

func connectWithSQL() *sql.DB {
	fmt.Println("=====Connect To Database=====")

	cfg := LoadDBConfig()
	dsn := cfg.ToDSN()
	fmt.Println("Connecting with DSN:", dsn)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		fmt.Println("❌ Failed to connect to database: %v", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		fmt.Println("❌ Database is not reachable: %v", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	SQLDB = db
	fmt.Println("✅ Successfully connected to database!")

	// Manual migration
	ResetDB(SQLDB)

	return SQLDB
}

func ResetDB(db *sql.DB) {
	fmt.Println("=== NATIVE MIGRATION ===")

	// Drop tables
	for _, name := range []string{"access_tokens", "refresh_tokens", "users"} {
		dropSQL := fmt.Sprintf(`DROP TABLE IF EXISTS "%s" CASCADE;`, name)
		if _, err := db.Exec(dropSQL); err != nil {
			fmt.Println("❌ Drop failed: %v", err)
		}
	}

	// Create tables
	createQueries := []string{
		GenerateCreateTableSQL("users", users.User{}),
		GenerateCreateTableSQL("access_tokens", auth.AccessToken{}),
		GenerateCreateTableSQL("refresh_tokens", auth.RefreshToken{}),
	}

	for _, q := range createQueries {

		if _, err := db.Exec(q); err != nil {
			fmt.Println("❌ Create failed: %v", err)
		}
	}

	fmt.Println("✅ Migrated using native SQL")
}
