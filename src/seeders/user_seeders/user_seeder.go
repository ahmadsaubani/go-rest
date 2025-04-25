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
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
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
		fmt.Println("❌ Batch insert failed: %w", err)
	}

	elapsed := time.Since(start).Seconds()
	fmt.Printf("✅ Seeded %d users in %.2f seconds\n", target-int64(userCount), elapsed)
}
