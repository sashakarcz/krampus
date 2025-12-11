package models

import (
	"time"
)

type User struct {
	ID           int64      `json:"id"`
	Username     string     `json:"username"`
	PasswordHash *string    `json:"-"` // Never send password hash to client
	Role         string     `json:"role"` // "ADMIN" or "USER"
	OIDCSubject  *string    `json:"oidc_subject,omitempty"`
	Email        *string    `json:"email,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	LastLogin    *time.Time `json:"last_login,omitempty"`
}

type UserRole string

const (
	RoleAdmin UserRole = "ADMIN"
	RoleUser  UserRole = "USER"
)

// IsAdmin checks if the user has admin role
func (u *User) IsAdmin() bool {
	return u.Role == string(RoleAdmin)
}
