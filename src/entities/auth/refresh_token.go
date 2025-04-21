package auth

import "time"

type RefreshToken struct {
	ID            uint      `gorm:"primaryKey"`
	UserID        uint      `gorm:"not null"`
	Token         string    `gorm:"not null;unique"`
	AccessTokenID uint      `gorm:"not null"` // Reference to AccessToken
	ExpiresAt     time.Time `gorm:"not null"`
	Claimed       bool      `gorm:"default:false"`
	CreatedAt     time.Time `gorm:"not null"`
	UpdatedAt     time.Time `gorm:"not null"`
}
