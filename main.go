package main

import (
	"gin/src/configs/database"
	"gin/src/routes"

	"github.com/gin-gonic/gin"
)

// main is the entry point of the application. It disables console colors for Gin,
// connects to the database, sets up the API routes, and starts the server on port 9000.

func main() {
	gin.DisableConsoleColor()

	// connection database
	database.ConnectDatabase()

	// routing to gin/src/routes folder
	r := routes.API()

	err := r.Run(":9000")
	if err != nil {
		panic(err)
	}
}
