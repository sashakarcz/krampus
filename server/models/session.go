package models

import (
	"time"
)

type Session struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	TokenHash string    `json:"-"` // Never send token hash to client
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}
