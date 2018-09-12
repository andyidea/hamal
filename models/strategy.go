package models

import "time"

type Strategy struct {
	ID        string     `gorm:"size:20;primary_key" json:"id"`
	Name      string     `gorm:"size:255;not null" json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}

type StrategyInstance struct {
	ID               uint `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	UserID           uint `gorm:"not null" json:"user_id"`
	User             User
	StrategyID       string `gorm:"size:20" json:"strategy_id"`
	Strategy         Strategy
	BaseCurrencyID   string `gorm:"size:20" json:"base_currency_id"`
	BaseCurrency     Currency
	TargetCurrencyID string `gorm:"size:20" json:"target_currency_id"`
	TargetCurrency   Currency
	Amount           float64    `gorm:"type:decimal(65,4)" json:"amount"`
	Interval         float64    `gorm:"type:decimal(65,4)" json:"interval"`
	ParamSeri        string     `gorm:"type:varchar(2550)" json:"param_seri"`
	UserExchangeID   uint       `json:"user_exchange_id"`
	Status           string     `gorm:"size:20" json:"status"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	DeletedAt        *time.Time `sql:"index" json:"deleted_at"`
}

type StrategyInstanceRecard struct {
	ID                 uint       `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	RecardType         string     `gorm:"size:20" json:"recard_type"`
	StrategyInstanceID uint       `gorm:"not null" json:"strategy_instance_id"`
	OrderID            string     `gorm:"size:20" json:"order_id"`
	Price              float64    `gorm:"type:decimal(65,2)" json:"price"`
	Amount             float64    `gorm:"type:decimal(65,4)" json:"amount"`
	Fee                float64    `gorm:"type:decimal(65,2)" json:"fee"`
	OrderTime          time.Time  `json:"order_time"`
	OrderStatus        int        `json:"order_status"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	DeletedAt          *time.Time `sql:"index" json:"deleted_at"`
}
