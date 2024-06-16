package models

import (
	"time"

	"gorm.io/gorm"
)

type Todo struct {
	ID        string `gorm:"primaryKey" json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
  CreatedAt time.Time `json:"created_at"`
  UpdatedAt time.Time`json:"updated_at"`
	Deleted   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}