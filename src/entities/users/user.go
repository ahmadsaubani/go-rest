package users

import "time"

type User struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" db:"id,primary,serial" json:"id"`
	Email     string    `gorm:"size:255;unique;not null" db:"email" json:"email" binding:"required,email"`
	Username  string    `gorm:"size:255;unique;not null" db:"username" json:"username" binding:"required,min=3,max=255"`
	Password  string    `gorm:"size:255;not null" db:"password" json:"password" binding:"required,min=6"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
type ResponseRegister struct {
	ID       int64  `db:"id" json:"id"`
	Email    string `db:"email" json:"email" binding:"required,email"`
	Username string `db:"username" json:"username" binding:"required,min=3,max=255"`
}

type ProfileResponse struct {
	ID       int64  `db:"id" json:"id"`
	Email    string `db:"email" json:"email" binding:"required,email"`
	Username string `db:"username" json:"username" binding:"required,min=3,max=255"`
}
