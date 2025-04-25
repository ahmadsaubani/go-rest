package auth

import "time"

type RefreshToken struct {
	UUID          string    `gorm:"uniqueIndex" db:"uuid" json:"uuid"`
	ID            int64     `gorm:"primaryKey" db:"id,primary,serial" json:"id"`
	UserID        int64     `gorm:"not null" db:"user_id"`
	Token         string    `gorm:"not null;unique" db:"token"`
	AccessTokenID int64     `gorm:"not null" db:"access_token_id"` // Reference to AccessToken
	ExpiresAt     time.Time `gorm:"not null" db:"expires_at"`
	Claimed       bool      `gorm:"default:false" db:"claimed"`
	CreatedAt     time.Time `gorm:"autoCreateTime" db:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" db:"updated_at"`
}
