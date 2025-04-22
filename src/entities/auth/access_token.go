package auth

import (
	"gin/src/entities/users"
	"time"
)

type AccessToken struct {
	ID        int64       `gorm:"primaryKey;autoIncrement" db:"id,primary,serial" json:"id"`
	UserID    int64       `gorm:"not null;index" db:"user_id" json:"user_id"`
	User      *users.User `gorm:"foreignKey:UserID db:"-"` // Tetap db:"-" karena User struct tidak bisa masuk langsung ke kolom SQL
	Token     string      `gorm:"uniqueIndex" db:"token" json:"token"`
	ExpiresAt time.Time   `db:"expires_at" json:"expires_at"`
	CreatedAt time.Time   `gorm:"autoCreateTime" db:"created_at" json:"created_at"`
	UpdatedAt time.Time   `gorm:"autoUpdateTime" db:"updated_at" json:"updated_at"`
}
