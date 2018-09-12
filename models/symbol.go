package models

import "time"

type Symbol struct {
	ID              string `gorm:"size:40;primary_key" json:"id"`
	Name            string `gorm:"size:255" json:"name"`
	BaseCurrencyID  string `gorm:"size:20" json:"base_currency_id"`
	BaseCurrency    Currency
	QuoteCurrencyID string `gorm:"size:20" json:"quote_currency_id"`
	QuoteCurrency   Currency
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `sql:"index" json:"deleted_at"`
}
