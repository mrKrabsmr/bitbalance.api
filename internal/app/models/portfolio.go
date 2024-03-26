package models

import (
	"time"

	"github.com/google/uuid"
)

type Portfolio struct {
	ID                   uuid.UUID `db:"id" json:"id"`
	UserID               uuid.UUID `db:"user_id" json:"user_id"`
	CMCCryptocurrencyID  int64     `db:"cmc_cryptocurrency_id" json:"cmc_cryptocurrency_id"`
	Cryptocurrency       string    `db:"cryptocurrency" json:"cryptocurrency"`
	CryptocurrencySymbol string    `db:"cryptocurrency_symbol" json:"cryptocurrency_symbol"`

	Price        float64   `db:"price" json:"price"`
	Count        float64   `db:"count" json:"count"`
	PurchaseTime time.Time `db:"purchase_time" json:"purchase_time"`
	Commentary   string    `db:"commentary" json:"commentary"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

func (p *Portfolio) TableName() string {
	return "portfolios"
}
