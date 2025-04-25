package user_seeders

import (
	"fmt"
	"gin/src/configs/database"
	"gin/src/entities/users"
	"gin/src/helpers"
	"log"
	"math/rand"
	"time"

	"github.com/bxcodec/faker/v4"
	"golang.org/x/crypto/bcrypt"
)

// func SeedUsers(userTargetCount int64) {
// 	userCount, err := helpers.CountModel[users.User]()
// 	if err != nil {
// 		fmt.Println("Error counting users:", err)
// 		return
// 	}

// 	if userCount >= userTargetCount {
// 		fmt.Printf("Users already seeded. Current count: %d\n", userCount)
// 		return
// 	}
// 	start := time.Now()
// 	// Seed users
// 	for i := userCount; i < userTargetCount; i++ {
// 		password := "password123"
// 		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

// 		email := faker.Email()
// 		if i == 0 {
// 			email = "ahmadsaubani@testing.com"
// 		}

// 		user := users.User{
// 			Email:    email,
// 			Username: fmt.Sprintf("%s_%d", faker.Username(), rand.Intn(10000)),
// 			Password: string(hashedPassword),
// 		}

// 		if err := helpers.InsertModel(&user); err != nil {
// 			log.Printf("failed to seed user %d: %v\n", i, err)
// 		}
// 	}

// 	// Calculate the total elapsed time after the loop finishes
// 	duration := time.Since(start)
// 	fmt.Printf("Total time to seed users: %s\n", duration.String())

// 	fmt.Printf("✅ Seeded %d users successfully\n", userTargetCount)
// }

func SeedUsers(db *database.DBConnection, target int64) {
	start := time.Now()

	userCount, err := helpers.CountModel[users.User]()
	if err != nil {
		log.Println("❌ Error counting users:", err)
		return
	}
	if int64(userCount) >= target {
		fmt.Println("✅ Users already seeded.")
		return
	}

	fmt.Printf("🔄 Seeding %d users...\n", target-int64(userCount))

	var usersBatch []users.User
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.MinCost)
	now := time.Now()

	for i := int64(userCount); i < target; i++ {
		email := faker.Email()
		if i == 0 {
			email = "ahmadsaubani@testing.com"
		}

		usersBatch = append(usersBatch, users.User{
			Email:     email,
			Username:  fmt.Sprintf("%s_%d", faker.Username(), rand.Intn(10000)),
			Password:  string(hashedPassword),
			CreatedAt: now,
			UpdatedAt: now,
		})
	}

	// Menggunakan fungsi InsertModelBatch untuk memasukkan data batch
	err = helpers.InsertModelBatch(usersBatch)
	if err != nil {
		log.Println("❌ Batch insert failed:", err)
	}

	elapsed := time.Since(start).Seconds()
	fmt.Printf("✅ Seeded %d users in %.2f seconds\n", target-int64(userCount), elapsed)
}
