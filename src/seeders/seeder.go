package seeders

import (
	"gin/src/configs/database"
	"gin/src/seeders/user_seeders"
)

func Run(db *database.DBConnection) {
	user_seeders.SeedUsers(db, 5000)
}
