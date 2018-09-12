package models

import "time"

type Currency struct {
	ID        string     `gorm:"size:20;primary_key" json:"id"`
	Name      string     `gorm:"size:255" json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}
