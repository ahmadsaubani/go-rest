package cmd

import (
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"

	"gin/src/seeders"
	"gin/src/utils/setups"
)

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed DB data",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load env
		if err := godotenv.Load(); err != nil {
			return err
		}

		// Connect DB tanpa migrasi
		dbConn := setups.DatabaseSetup(false, true)

		// Jalankan seeder secara terpisah
		seeders.Run(dbConn)
		return nil
	},
}
