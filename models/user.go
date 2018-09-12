package models

import "time"

//用户表
type User struct {
	ID          uint       `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Email       string     `gorm:"size:255" json:"email"`
	Cellphone   string     `gorm:"size:255" json:"cellphone"`
	Username    string     `gorm:"size:100;not null;unique" json:"username"`
	Name        string     `gorm:"size:255" json:"name"`
	PasswordMd5 string     `gorm:"size:255;not null;" json:"password_md5"`
	Avatar      string     `gorm:"size:255" json:"avatar"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `sql:"index" json:"deleted_at"`
}

//用户交易所表
type UserExchange struct {
	ID        uint       `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	UserID    uint       `json:"user_id"`
	Exchange  string     `gorm:"size:255" json:"name"`
	Tag       string     `gorm:"size:255" json:"tag"`
	AccessKey string     `gorm:"size:255" json:"access_key"`
	SecretKey string     `gorm:"size:255" json:"secret_key"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}
