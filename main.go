package main

import (
	"gin/src/configs/database"
	"gin/src/configs/registrations"
	"gin/src/routes"
	"gin/src/seeders"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// main is the entry point of the application. It disables console colors for Gin,
// connects to the database, sets up the API routes, and starts the server on port 9000.

func main() {
	// Disable console color for clean output
	gin.DisableConsoleColor()
	ginEngine := gin.Default()

	// registration global middleware
	ginEngine = registrations.GlobalMiddlewares(ginEngine)

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file: " + err.Error()) // Panic with the error message if .env file loading fails
	}

	// Establish database connection
	db := database.ConnectDatabase()

	// run seeders
	seeders.Run(db)

	// Initialize routes
	r := routes.API(db, ginEngine)

	// Run the server on port 9000
	if err := r.Run(":9000"); err != nil {
		panic("Error starting server: " + err.Error()) // Panic if server fails to start
	}
}

// loggers.InitLogger()

// // Menulis log
// loggers.Log.Info("App started")
// loggers.Log.Warn("This is a warning")
// loggers.Log.Error("Something went wrong")

// log := loggers.NewLogger()
// defer log.Close()

// debug exit
// helpers.DdLog("Debugging user before login:", ser)
