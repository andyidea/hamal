package models

import "time"

type Exchange struct {
	Name      string     `gorm:"size:255;primary_key" json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}
