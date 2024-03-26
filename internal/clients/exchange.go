package clients

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	core "fl/my-portfolio/internal/app"
	"fl/my-portfolio/internal/configs"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type ExchangeClient struct {
	config configs.ExchangeConfig
}

func NewExchangeClient() *ExchangeClient {
	return &ExchangeClient{
		config: core.GetConfig().Exchange,
	}
}

func (c *ExchangeClient) generateSignature(message, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(message))
	signingKey := fmt.Sprintf("%x", mac.Sum(nil))
	return signingKey
}

func (c *ExchangeClient) generateSignatureOkx(message, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(message))
	signingKey := fmt.Sprintf("%x", mac.Sum(nil))
	signingKey = base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return signingKey
}

func (c *ExchangeClient) GetBinanceStakingData() ([]StakingData, error) {
	apiKey := c.config.BinanceAPIKey
	secretKey := c.config.BinanceSecretKey

	timestamp := time.Now().Unix() * 1000
	baseURL := "https://api.binance.com"
	endpoint := "/sapi/v1/staking/productList"
	method := "GET"

	params := "?timestamp=" + fmt.Sprintf("%d", timestamp) + "&product=STAKING"
	signature := c.generateSignature(params[1:], secretKey)
	params += "&signature=" + signature
	req, err := http.NewRequest(method, baseURL+endpoint+params, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-MBX-APIKEY", apiKey)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var objects []*BinanceStakingDataResponse

	if err = json.Unmarshal(body, &objects); err != nil {
		return nil, err
	}

	stakingDataList := make([]StakingData, 0, len(objects))
	for _, obj := range objects {
		apy, err := strconv.ParseFloat(obj.Detail.Apy, 64)
		if err != nil {
			return nil, err
		}

		minimum, err := strconv.ParseFloat(obj.Quota.Minimum, 64)
		if err != nil {
			return nil, err
		}

		stakingData := StakingData{
			Asset:        obj.Detail.Asset,
			RewardAsset:  obj.Detail.RewardAsset,
			Duration:     obj.Detail.Duration,
			Apy:          apy,
			QuotaMinimum: minimum,
		}

		stakingDataList = append(stakingDataList, stakingData)
	}

	return stakingDataList, nil
}

// func (c *ExchangeClient) GetBybitStakingData() {
// 	timestamp := time.Now().Unix() * 1000
// 	baseURL := "https://api.bybit.com"
// 	endpoint := "/v5/broker/earnings-info"
// 	method := "GET"

// 	params := "?timestamp=" + fmt.Sprintf("%d", timestamp) + "&api_key=" + apiKey
// 	queryString := "api_key=" + apiKey + "&timestamp=" + fmt.Sprintf("%d", timestamp)
// 	signature := c.generateSignature(queryString, secretKey)
// 	params += "&sign=" + signature
// 	req, err := http.NewRequest(method, baseURL+endpoint+params, nil)
// 	if err != nil {
// 		fmt.Println("Error creating request:", err)
// 		return
// 	}

// 	req.Header.Add("X-MBX-APIKEY", apiKey)
// 	req.Header.Add("X-BAPI-SIGN", signature)
// 	req.Header.Add("X-BAPI-TIMESTAMP", fmt.Sprintf("%d", timestamp))
// 	req.Header.Add("Content-Type", "application/json")

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		fmt.Println("Error sending request:", err)
// 		return
// 	}
// 	defer resp.Body.Close()

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		fmt.Println("Error reading response:", err)
// 		return
// 	}

// 	fmt.Println(resp.StatusCode)

// 	fmt.Println(string(body))
// }

func (c *ExchangeClient) GetOkxStakingData() ([]StakingData, error) {
	apiKey := c.config.OkxAPIKey
	secretKey := c.config.OkxSecretKey
	passphrase := c.config.OkxPassphrase

	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	baseURL := "https://www.okx.com"
	endpoint := "/api/v5/finance/staking-defi/offers"
	method := "GET"

	queryString := timestamp + method + endpoint
	fmt.Println(queryString)
	signature := c.generateSignatureOkx(queryString, secretKey)
	req, err := http.NewRequest(method, baseURL+endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("OK-ACCESS-KEY", apiKey)
	req.Header.Add("OK-ACCESS-SIGN", signature)
	req.Header.Add("OK-ACCESS-TIMESTAMP", timestamp)
	req.Header.Add("OK-ACCESS-PASSPHRASE", passphrase)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var objects *OkxStakingDataResponse

	if err = json.Unmarshal(body, &objects); err != nil {
		return nil, err
	}

	stakingDataList := make([]StakingData, 0, len(objects.Data))
	for _, obj := range objects.Data {
		for _, subobjAsset := range obj.InvestData {
			duration, err := strconv.Atoi(obj.Term)
			if err != nil {
				return nil, err
			}

			minimum, err := strconv.ParseFloat(subobjAsset.MinAmt, 64)
			if err != nil {
				return nil, err
			}

			apy, err := strconv.ParseFloat(obj.Apy, 64)
			if err != nil {
				return nil, err
			}

			for _, subobjRewardAsset := range obj.EarningData {
				stakingData := StakingData{
					Asset:        subobjAsset.Ccy,
					RewardAsset:  subobjRewardAsset.Ccy,
					Duration:     duration,
					Apy:          apy,
					QuotaMinimum: minimum,
				}

				stakingDataList = append(stakingDataList, stakingData)
			}
		}
	}

	return stakingDataList, nil
}
