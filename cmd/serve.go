package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"gin/src/configs/database"
	"gin/src/configs/registrations"
	"gin/src/routes"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run HTTP server",
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = godotenv.Load()
		dbConn := database.ConnectDatabase(false)
		ginEngine := gin.Default()
		ginEngine = registrations.GlobalMiddlewares(ginEngine)
		r := routes.API(dbConn, ginEngine)
		return r.Run(":9000")
	},
}
