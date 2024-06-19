package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        string         `gorm:"primaryKey" json:"id"`
	LoginID   string         `gorm:"uniqueIndex;not null" json:"login_id"`
	Name      string         `gorm:"not null" json:"username"`
	Password  string         `json:"-"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	Deleted   gorm.DeletedAt `gorm:"index" json:"-"`
}

type UserDTO struct {
	LoginID  string `json:"login_id"`
	Username string `json:"username"`
	Password string `json:"password"`
}
