package user_seeders

import (
	"fmt"
	"gin/src/entities/users"
	"gin/src/helpers"
	"log"
	"math/rand"

	"github.com/bxcodec/faker/v4"
	"golang.org/x/crypto/bcrypt"
)

func SeedUsers() {
	userCount, err := helpers.CountModel[users.User]()
	if err != nil {
		fmt.Println("Error counting users:", err)
		return
	}

	if userCount >= 1000 {
		fmt.Println("Users already seeded.")
		return
	}

	for i := 0; i < 10; i++ {
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

		helpers.DumpLog("creating user:", user)

		if err := helpers.InsertModel(&user); err != nil {
			log.Printf("failed to seed user %d: %v\n", i, err)
		}
	}

	fmt.Println("✅ Seeded 1000 users successfully")
}
