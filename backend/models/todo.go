package models

import (
	"gorm.io/gorm"
)

type Todo struct {
	ID        string `gorm:"primaryKey" json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
	Deleted   gorm.DeletedAt `json:"deleted_at,omitempty"`
}