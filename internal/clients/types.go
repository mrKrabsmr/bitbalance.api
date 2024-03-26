package clients

// api endpoints
var (
	cryptoListingLatest = "https://pro-api.coinmarketcap.com/v1/cryptocurrency/listings/latest"
)

type CryptocurrencyDataDetail struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
	Quote  struct {
		USD struct {
			Price float64 `json:"price"`
			// Volume24h        float64 `json:"volume_24h"`
			// VolumeChange24h  float64 `json:"volume_change_24h"`
			// PercentChange1h  float64 `json:"percent_change_1h"`
			PercentChange24h float64 `json:"percent_change_24h"`
			// PercentChange7d  float64 `json:"percent_change_7d"`
			PercentChange30d float64 `json:"percent_change_30d"`
			// PercentChange60d float64 `json:"percent_change_60d"`
			PercentChange90d float64 `json:"percent_change_90d"`
		} `json:"USD"`
		BTC struct {
			Price float64 `json:"price"`
			// Volume24h        float64 `json:"volume_24h"`
			// VolumeChange24h  float64 `json:"volume_change_24h"`
			// PercentChange1h  float64 `json:"percent_change_1h"`
			PercentChange24h float64 `json:"percent_change_24h"`
			// PercentChange7d  float64 `json:"percent_change_7d"`
			PercentChange30d float64 `json:"percent_change_30d"`
			// PercentChange60d float64 `json:"percent_change_60d"`
			PercentChange90d float64 `json:"percent_change_90d"`
		} `json:"BTC"`
		RUB struct {
			Price float64 `json:"price"`
			// Volume24h        float64 `json:"volume_24h"`
			// VolumeChange24h  float64 `json:"volume_change_24h"`
			// PercentChange1h  float64 `json:"percent_change_1h"`
			PercentChange24h float64 `json:"percent_change_24h"`
			// PercentChange7d  float64 `json:"percent_change_7d"`
			PercentChange30d float64 `json:"percent_change_30d"`
			// PercentChange60d float64 `json:"percent_change_60d"`
			PercentChange90d float64 `json:"percent_change_90d"`
		} `json:"RUB"`
	} `json:"quote"`
}

type CryptoListingsLatestResponse struct {
	Status struct {
		TotalCount int `json:"total_count"`
	} `json:"status"`
	Data []CryptocurrencyDataDetail `json:"data"`
}

type BinanceStakingDataResponse struct {
	Detail struct {
		Asset       string `json:"asset"`
		RewardAsset string `json:"rewardAsset"`
		Duration    int    `json:"duration"`
		Apy         string `json:"apy"`
	} `json:"detail"`
	Quota struct {
		TotalPersonalQuota string `json:"totalPersonalQuota"`
		Minimum            string `json:"minimum"`
	} `json:"quota"`
}

type OkxStakingDataResponse struct {
	Data []struct {
		Ccy        string `json:"ccy"`
		Term       string `json:"term"`
		Apy        string `json:"apy"`
		InvestData []struct {
			Ccy    string `json:"ccy"`
			MinAmt string `json:"minAmt"`
		} `json:"investData"`
		EarningData []struct {
			Ccy string `json:"ccy"`
		} `json:"earningData`
	} `json:"data"`
}

type StakingData struct {
	Asset        string  `json:"asset"`
	RewardAsset  string  `json:"reward_asset"`
	Duration     int     `json:"duration"`
	Apy          float64 `json:"apy"`
	QuotaMinimum float64 `json:"quota_minimum"`
}