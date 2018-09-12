package models

import "time"

type DailyStat struct {
	UserID      uint       `gorm:"primary_key" json:"user_id"`
	Date        time.Time  `gorm:"type:date;primary_key" json:"date"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `sql:"index" json:"deleted_at"`
	BTCBalance  float64    `gorm:"type:decimal(65,4)" json:"btc_balance"`
	USDTBalance float64    `gorm:"type:decimal(65,4)" json:"usdt_balance"`
	Profit      float64    `gorm:"type:decimal(65,2)" json:"profit"`
}
