package auth

import (
	"gin/src/entities/users"
	"time"
)

type AccessToken struct {
	UUID      string      `gorm:"uniqueIndex" db:"uuid" json:"uuid"`
	ID        int64       `gorm:"primaryKey;autoIncrement" db:"id,primary,serial" json:"id"`
	UserID    int64       `gorm:"not null;index" db:"user_id" json:"user_id"`
	User      *users.User `gorm:"foreignKey:UserID" db:"-" json:"user"`
	Token     string      `gorm:"uniqueIndex" db:"token" json:"token"`
	ExpiresAt time.Time   `db:"expires_at" json:"expires_at"`
	Revoked   bool        `gorm:"default:false" db:"revoked" json:"revoked"`
	CreatedAt time.Time   `gorm:"autoCreateTime" db:"created_at" json:"created_at"`
	UpdatedAt time.Time   `gorm:"autoUpdateTime" db:"updated_at" json:"updated_at"`
}
