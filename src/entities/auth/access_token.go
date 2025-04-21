package auth

import (
	"gin/src/entities/users"
	"time"
)

type AccessToken struct {
	ID        uint       `gorm:"primaryKey;autoIncrement"`
	UserID    uint       `gorm:"not null;index"`
	User      users.User `gorm:"foreignKey:UserID"`
	Token     string     `gorm:"uniqueIndex"`
	ExpiresAt time.Time
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
