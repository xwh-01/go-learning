package models

import (
	"gorm.io/gorm"
	"time"
)

type Exchangerate struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	CreatedAt    time.Time      `gorm:"index" json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	FromCurrency string         `gorm:"type:varchar(10);not null;index:idx_currency,priority:1" json:"from_currency"`
	ToCurrency   string         `gorm:"type:varchar(10);not null;index:idx_currency,priority:2" json:"to_currency"`
	Rate         float64        `gorm:"type:decimal(20,8);not null" json:"rate"`
	Date         time.Time      `gorm:"index" json:"date"`
}

func (Exchangerate) TableName() string {
	return "exchange_rates"
}
