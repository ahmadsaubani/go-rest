package main

import (
	"gin/src/configs/database"
	"gin/src/routes"
	"gin/src/seeders/user_seeders"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// main is the entry point of the application. It disables console colors for Gin,
// connects to the database, sets up the API routes, and starts the server on port 9000.

func main() {
	gin.DisableConsoleColor()

	// connection database
	db := database.ConnectDatabase()

	// seeder
	user_seeders.SeedUsers()

	// routing to gin/src/routes folder
	r := routes.API(db)

	err := r.Run(":9000")
	if err != nil {
		panic(err)
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
