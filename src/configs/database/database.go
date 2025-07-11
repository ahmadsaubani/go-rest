package database

import (
	"database/sql"
	"fmt"
	"gin/src/entities"
	utils "gin/src/utils/tablers"
	"os"
	"time"

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

func ConnectDatabase(runMigration bool) *DBConnection {
	fmt.Println("===== Connecting To Database =====")

	if os.Getenv("USE_GORM") == "true" {
		gdb := connectWithGORM()
		if runMigration {
			ensureMigrations(true, gdb, nil)
		}
		conn := &DBConnection{Gorm: gdb}

		return conn
	}

	sdb := connectWithSQL()
	if runMigration {
		ensureMigrations(false, nil, sdb)
	}
	conn := &DBConnection{SQL: sdb}

	return conn
}

// connectWithGORM hanya membuka koneksi GORM, tanpa drop/recreate.
func connectWithGORM() *gorm.DB {
	fmt.Println("=====USING GORM=====")
	cfg := LoadDBConfig()
	dsn := cfg.ToDSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("❌ Failed to connect to database using GORM: %w", err)
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Hour)
	GormDB = db
	fmt.Println("✅ GORM connected")
	return db
}

// connectWithSQL hanya membuka koneksi database/sql
func connectWithSQL() *sql.DB {
	fmt.Println("=====USING NATIVE SQL=====")
	cfg := LoadDBConfig()
	dsn := cfg.ToDSN()

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		fmt.Println("❌ Failed to connect to database using Native SQL: %w", err)
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)
	if err := db.Ping(); err != nil {
		fmt.Println("❌ ping failed: %w", err)

	}
	SQLDB = db
	fmt.Println("✅ database/sql connected")
	return db
}

func ensureMigrations(useGORM bool, gdb *gorm.DB, sdb *sql.DB) {
	gormModels := make([]interface{}, 0, len(entities.RegisteredEntities))
	sqlTables := make([]string, 0, len(entities.RegisteredEntities))

	for _, model := range entities.RegisteredEntities {
		gormModels = append(gormModels, model)
		sqlTables = append(sqlTables, utils.GetTableName(model))
	}

	if useGORM {
		migrator := gdb.Migrator()

		allExist := true
		for _, m := range gormModels {
			if !migrator.HasTable(m) {
				allExist = false
				break
			}
		}
		if allExist {
			fmt.Println("ℹ️  All tables already exist (GORM), skipping migration")
			return
		}

		fmt.Println("🔧 Running GORM AutoMigrate…")
		if err := gdb.AutoMigrate(gormModels...); err != nil {
			fmt.Printf("❌ GORM AutoMigrate failed: %v\n", err)
		} else {
			fmt.Println("✅ GORM migration complete")
		}
	} else {
		allExist := true
		for _, tbl := range sqlTables {
			var count int
			err := sdb.QueryRow(`
				SELECT COUNT(*) 
				  FROM information_schema.tables 
				 WHERE table_schema = 'public' AND table_name = $1
			`, tbl).Scan(&count)
			if err != nil || count == 0 {
				allExist = false
				break
			}
		}
		if allExist {
			fmt.Println("ℹ️  All tables already exist (native), skipping migration")
			return
		}

		fmt.Println("🔧 Running native SQL migrations…")

		// DROP
		for _, name := range sqlTables {
			_, _ = sdb.Exec(fmt.Sprintf(`DROP TABLE IF EXISTS "%s" CASCADE;`, name))
		}

		// CREATE
		for i, model := range entities.RegisteredEntities {
			query := GenerateCreateTableSQL(sqlTables[i], model)
			if _, err := sdb.Exec(query); err != nil {
				fmt.Printf("❌ failed to exec create table [%s]: %v\n", sqlTables[i], err)
			}
		}
		fmt.Println("✅ Native SQL migration complete")
	}
}
