package controllers

import (
	"fl/my-portfolio/internal/clients"
	"fl/my-portfolio/pkg/types"

	"github.com/google/uuid"
)

var text = map[uint8]string{
	0: "server error, please wait for fixes",
	1: "invalid input",
	2: "authentication failed",
}

type ResponseSuccess struct {
	Error bool        `json:"error"`
	Data  interface{} `json:"data"`
}

type ResponseError struct {
	Error   bool   `json:"error"`
	Message string `json:"error_string"`
}

type Tokens struct {
	UserID       string `json:"user_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RegisterData struct {
	Email     string     `json:"email" validate:"required,email"`
	Password  string     `json:"password" validate:"required,min=8,max=20,alphanum"`
	FirstName string     `json:"first_name" validate:"required,min=2"`
	LastName  string     `json:"last_name" validate:"required,min=2"`
	Gender    string     `json:"gender" validate:"required,oneof=male female"`
	BirthDate types.Time `json:"birth_date" validate:"required,lt"`
}

type LoginData struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type RefreshData struct {
	RefreshToken string `json:"refresh_token" validate:"required,jwt"`
}

type PortfolioPostData struct {
	CMCCryptocurrencyID int64      `json:"cmc_cryptocurrency_id" validate:"required"`
	Price               float64    `json:"price" validate:"required"`
	Count               float64    `json:"count" validate:"required"`
	PurchaseTime        types.Time `json:"purchase_time" validate:"required"`
	Commentary          string     `json:"commentary"`
}

type PortfolioPatchData struct {
	Price        float64    `json:"price"`
	Count        float64    `json:"count"`
	PurchaseTime types.Time `json:"purchase_time"`
	Commentary   string     `json:"commentary"`
}

type PortfolioDetailGetResponse struct {
	ID                   uuid.UUID  `json:"id"`
	CryptID              int64      `json:"cryptocurrency_id"`
	Cryptocurrency       string     `json:"cryptocurrency"`
	CryptocurrencySymbol string     `json:"cryptocurrency_symbol"`
	Price                float64    `json:"price" validate:"required"`
	Count                float64    `json:"count" validate:"required"`
	Sum                  float64    `json:"sum"`
	NowPrice             float64    `json:"now_price"`
	NowSum               float64    `json:"now_sum"`
	PurchaseTime         types.Time `json:"purchase_time" validate:"required"`
	Commentary           string     `json:"commentary"`
	CreatedAt            types.Time `json:"created_at"`
	PercentChange24h     float64    `json:"percent_change_24h"`
	PercentChange30d     float64    `json:"percent_change_30d"`
	PercentChange90d     float64    `json:"percent_change_90d"`
	PortfolioShare       float64    `json:"portfolio_share"`
	ROI                  float64    `json:"ROI"`
}

type PortfolioGetResponse struct {
	TotalPurchaseSum float64                      `json:"total_purchased_sum"`
	TotalNowSum      float64                      `json:"total_now_sum"`
	ROI              float64                      `json:"ROI"`
	Cryptocurrencies []PortfolioDetailGetResponse `json:"cryptocurrencies"`
}

type StakingGetResponse struct {
	CrpyptocurrencyData clients.CryptocurrencyDataDetail `json:"cryptocurrency_data"`
	BinanceStaking      *clients.StakingData             `json:"binance_staking"`
	BybitStaking        *clients.StakingData             `json:"bybit_staking"`
	OkxStaking          *clients.StakingData             `json:"okx_staking"`
}

type StakingDetailGetResponse struct {
	BinanceStaking      []clients.StakingData             `json:"binance_stakings"`
	BybitStaking        []clients.StakingData             `json:"bybit_stakings"`
	OkxStaking          []clients.StakingData             `json:"okx_stakings"`
}
