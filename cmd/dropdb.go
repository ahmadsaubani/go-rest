package cmd

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"

	"gin/src/configs/database"
	"gin/src/entities"
	utils "gin/src/utils/tablers"
)

var dropDbCmd = &cobra.Command{
	Use:   "dropdb",
	Short: "Drop all tables (native SQL mode)",
	Run: func(_ *cobra.Command, _ []string) {
		// Load environment variable terlebih dahulu
		if err := godotenv.Load(); err != nil {
			fmt.Println("❌ Failed to load .env:", err)
			os.Exit(1)
		}

		// Paksa gunakan native SQL mode
		os.Setenv("USE_GORM", "false")

		// Connect ke database
		conn := database.ConnectDatabase(false)
		sqlDB := conn.SQL
		if sqlDB == nil {
			fmt.Println("❌ Native SQL DB connection not available")
			os.Exit(1)
		}

		// Loop dari RegisteredModels
		for _, model := range entities.RegisteredEntities {
			name := utils.GetTableNameRuntime(model)
			_, err := sqlDB.Exec(fmt.Sprintf(`DROP TABLE IF EXISTS "%s" CASCADE;`, name))
			if err != nil {
				fmt.Printf("❌ Failed to drop %s: %v\n", name, err)
			} else {
				fmt.Printf("✅ Dropped table: %s\n", name)
			}
		}
	},
}
