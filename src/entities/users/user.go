package users

import "time"

type User struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Email     string    `gorm:"size:255;unique;not null" json:"email" binding:"required,email"`
	Username  string    `gorm:"size:255;unique;not null" json:"username" binding:"required,min=3,max=255"`
	Password  string    `gorm:"size:255;not null" json:"password" binding:"required,min=6"` // Exclude password from JSON response
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type ResponseRegister struct {
	ID       int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Email    string `gorm:"size:255;unique;not null" json:"email" binding:"required,email"`
	Username string `gorm:"size:255;unique;not null" json:"username" binding:"required,min=3,max=255"`
}

type ProfileResponse struct {
	ID       int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Email    string `gorm:"size:255;unique;not null" json:"email" binding:"required,email"`
	Username string `gorm:"size:255;unique;not null" json:"username" binding:"required,min=3,max=255"`
}
