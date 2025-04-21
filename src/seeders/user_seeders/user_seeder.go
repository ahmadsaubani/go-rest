package user_seeders

import (
	"fmt"
	"gin/src/configs/database"
	"gin/src/entities/users"
	"gin/src/helpers"
	"log"
	"math/rand"

	"github.com/bxcodec/faker/v4"
	"golang.org/x/crypto/bcrypt"
)

func SeedUsers() {
	var userCount int64
	database.DB.Model(&users.User{}).Count(&userCount)

	if userCount >= 1000 {
		fmt.Println("Users already seeded.")
		return
	}

	for i := 0; i < 1000; i++ {
		password := "password123" // default password
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		email := faker.Email()
		if i == 0 {
			email = "ahmadsaubani@testing.com"
		}

		user := users.User{
			Email:    email,
			Username: fmt.Sprintf("%s_%d", faker.Username(), rand.Intn(10000)),
			Password: string(hashedPassword),
		}

		helpers.DumpLog("User trying to login:", user)

		if err := database.DB.Create(&user).Error; err != nil {
			log.Printf("failed to seed user %d: %v\n", i, err)
		}
	}

	fmt.Println("✅ Seeded 1000 users successfully")
}
