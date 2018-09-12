package models

import "time"

type OrderHistory struct {
	ID             uint       `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	SourceID       uint       `json:"source_id"`
	ExchangeID     uint       `json:"exchange_id"`
	UserExchangeID uint       `json:"user_exchange_id"`
	CoinSymbol     string     `gorm:"size:255" json:"coin_symbol"`
	CurrencySymbol string     `gorm:"size:255"json:"currency_symbol"`
	Price          string     `gorm:"size:255"json:"price"`
	Amount         string     `gorm:"size:255"json:"amount"`
	Fee            string     `gorm:"size:255"json:"fee"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `sql:"index" json:"deleted_at"`
}
