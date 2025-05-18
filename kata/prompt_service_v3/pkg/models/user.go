package models

import (
	"time"
)

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"unique;not null"`
	Password  string    `json:"password" gorm:"not null"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Role      string    `json:"role" gorm:"default:'user'"` // e.g., 'admin', 'user'
	ExpiredAt time.Time `json:"expired_at"`                // expiration time for account
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}