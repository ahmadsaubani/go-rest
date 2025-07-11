package cmd

import (
	"github.com/spf13/cobra"
)

// rootCmd adalah entry point untuk CLI
var rootCmd = &cobra.Command{
	Use:   "apps",
	Short: "CLI untuk mengelola",
}

// Execute menjalankan root command
func Execute() error {
	return rootCmd.Execute()
}

// Inisialisasi subcommands
func init() {
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(seedCmd)
	rootCmd.AddCommand(dropDbCmd)
}
