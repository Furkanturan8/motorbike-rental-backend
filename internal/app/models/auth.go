package models

import (
	"github.com/google/uuid"
	_ "gorm.io/gorm"
	"time"
)

type AuthRefreshToken struct {
	TokenID   uuid.UUID `gorm:"type:char(36);primary_key" json:"token_id"`
	UserID    int64     `gorm:"index;not null" json:"user_id"`
	Role      float64   `gorm:"index;not null" json:"role"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
}

func (AuthRefreshToken) ModelName() string {
	return "user_refresh_token"
}

// AuthTokenPair defines the structure for access and refresh tokens
type AuthTokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
