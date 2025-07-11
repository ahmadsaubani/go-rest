package cmd

import (
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"

	"gin/src/configs/database"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run DB migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load env
		if err := godotenv.Load(); err != nil {
			return err
		}

		// Jalankan koneksi + migrasi (tanpa seeding)
		_ = database.ConnectDatabase(true) // hanya migrasi
		return nil
	},
}
