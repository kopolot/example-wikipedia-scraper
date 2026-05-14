package model

import (
	types "example-wikipedia-scraper/internal/types/model"
	"time"
)

const (
	RoleUser  types.UserRole = "user"
	RoleAdmin types.UserRole = "admin"
)

type User struct {
	Model
	LastLoginAt                 time.Time      `json:"lastLoginAt" gorm:"column:last_login_at"`
	PasswordResetTokenExpiresAt *time.Time     `json:"-" gorm:"column:password_reset_token_expires_at;default:null"`
	LastLogoutAt                *time.Time     `json:"lastLogoutAt" gorm:"column:last_logout_at;default:null"`
	Email                       string         `json:"email" gorm:"not null;uniqueIndex"`
	Password                    string         `json:"-" gorm:"not null"`
	Role                        types.UserRole `json:"role" gorm:"type:enum('user','admin');default:'user'"`
	Username                    string         `json:"username" gorm:"not null;uniqueIndex"`
	EmailVerificationToken      string         `json:"-" gorm:"column:email_verification_token"`
	PasswordResetToken          string         `json:"-" gorm:"column:password_reset_token"`
	EmailVerified               bool           `json:"emailVerified" gorm:"column:email_verified;default:false"`
}
