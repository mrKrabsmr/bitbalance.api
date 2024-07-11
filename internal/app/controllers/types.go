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
	Username     string `json:"username"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RegisterData struct {
	Username string `json:"username" validate:"required,lowercase"`
	Password string `json:"password" validate:"required,min=8,max=20,alphanum"`
}

type LoginData struct {
	Username string `json:"username" validate:"required"`
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

type PurchaseGetResponse struct {
	ID           uuid.UUID  `json:"id"`
	Price        float64    `json:"price" validate:"required"`
	Count        float64    `json:"count" validate:"required"`
	Sum          float64    `json:"sum"`
	PurchaseTime types.Time `json:"purchase_time" validate:"required"`
	Commentary   string     `json:"commentary"`
	CreatedAt    types.Time `json:"created_at"`
}

type PortfolioDetailGetResponse struct {
	CryptID              int64                 `json:"cryptocurrency_id"`
	Cryptocurrency       string                `json:"cryptocurrency"`
	CryptocurrencySymbol string                `json:"cryptocurrency_symbol"`
	NowPrice             float64               `json:"now_price"`
	Count                float64               `json:"count"`
	Sum                  float64               `json:"sum"`
	PercentChange24h     float64               `json:"percent_change_24h"`
	PercentChange30d     float64               `json:"percent_change_30d"`
	PercentChange90d     float64               `json:"percent_change_90d"`
	PortfolioShare       float64               `json:"portfolio_share"`
	ROI                  float64               `json:"ROI"`
	Profit               float64               `json:"profit"`
	Purchases            []PurchaseGetResponse `json:"purchases"`
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
	BinanceStaking []clients.StakingData `json:"binance_stakings"`
	BybitStaking   []clients.StakingData `json:"bybit_stakings"`
	OkxStaking     []clients.StakingData `json:"okx_stakings"`
}
