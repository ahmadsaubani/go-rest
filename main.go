package main

import (
	"fmt"
	"os"

	"gin/cmd"

	_ "github.com/lib/pq"
)

// main is the entry point of the application. It disables console colors for Gin,
/*************  ✨ Windsurf Command ⭐  *************/
// main is the entry point of the application. It disables console colors for Gin,
// connects to the database, runs seeders, sets up the API routes, and starts the server on port 9000.
//
// It first executes the root command from the subcommand package to parse any command-line arguments.
// It then sets up Gin's default engine, registers global middleware, loads environment variables from
// the .env file, and establishes a connection to the database. It runs the seeders package to populate
// the database with initial data. Finally, it sets up the API routes and starts the server on port 9000.
/*******  779356ef-c832-4676-b360-9c30fcb57db5  *******/ // connects to the database, sets up the API routes, and starts the server on port 9000.

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// // Disable console color for clean output
	// gin.DisableConsoleColor()
	// ginEngine := gin.Default()

	// // registration global middleware
	// ginEngine = registrations.GlobalMiddlewares(ginEngine)

	// // Load environment variables from .env file
	// if err := godotenv.Load(); err != nil {
	// 	panic("Error loading .env file: " + err.Error()) // Panic with the error message if .env file loading fails
	// }

	// // Establish database connection
	// db := database.ConnectDatabase()

	// // run seeders
	// seeders.Run(db)

	// // Initialize routes
	// r := routes.API(db, ginEngine)

	// // Run the server on port 9000
	// if err := r.Run(":9000"); err != nil {
	// 	panic("Error starting server: " + err.Error()) // Panic if server fails to start
	// }
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
