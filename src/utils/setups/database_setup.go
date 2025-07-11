package setups

import (
	"gin/src/configs/database"
	"gin/src/seeders"
)

func DatabaseSetup(runMigrate, runSeed bool) *database.DBConnection {
	conn := database.ConnectDatabase(runMigrate)
	if runSeed {
		seeders.Run(conn)
	}
	return conn
}
