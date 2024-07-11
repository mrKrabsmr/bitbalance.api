package services

import (
	"context"
	"encoding/json"
	"errors"
	"fl/my-portfolio/internal/app/models"
	"fl/my-portfolio/internal/clients"
	"sort"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

func (s *Service) GetAllCryptocurrencies(isSorted bool) (*clients.CryptoListingsLatestResponse, error) {
	var arg string
	if isSorted {
		arg = "cryptocurrencies_sorted"
	} else {
		arg = "cryptocurrencies"
	}

	cryptocurrencies, err := s.cacher.Get(context.Background(), arg).Result()
	if err != nil {
		s.logger.Info(err)
		if errors.Is(err, redis.Nil) {
			cryptocurrencies, err := s.cmcClient.GetCryptoListingsLatest()
			if err != nil {
				return nil, err
			}

			jsonData, err := json.Marshal(cryptocurrencies)
			if err != nil {
				return nil, err
			}
			data := cryptocurrencies.Data
			sort.Slice(data, func(i, j int) bool {
				return data[i].ID < data[j].ID
			})

			cryptocurrenciesSorted := cryptocurrencies
			cryptocurrenciesSorted.Data = data

			jsonDataSorted, _ := json.Marshal(cryptocurrenciesSorted)

			go func() {
				s.cacher.Set(context.Background(), "cryptocurrencies", jsonData, time.Minute*5)
				s.cacher.Set(context.Background(), "cryptocurrencies_sorted", jsonDataSorted, time.Minute*5)
			}()

			if isSorted {
				return cryptocurrenciesSorted, nil
			}

			return cryptocurrencies, nil
		} else {
			return nil, err
		}
	}

	var obj *clients.CryptoListingsLatestResponse

	jsonString, err := json.Marshal(cryptocurrencies)
	if err != nil {
		return nil, err
	}

	ss, _ := strconv.Unquote(string(jsonString))
	if err = json.Unmarshal([]byte(ss), &obj); err != nil {
		return nil, err
	}

	// s.logger.Info("Hello")
	// file, err := os.Create("cc.txt")
	// if err != nil {
	// 	panic(err)
	// }

	// defer file.Close()
	// for _, c := range obj.Data {
	// 	file.Write([]byte(strings.ToLower(c.Symbol) + "\n"))
	// }
	return obj, nil
}

func (s *Service) GetCryptocurrencyAdditionalData(
	portfolioDetail *models.Portfolio, totalNowSum float64, crypts *clients.CryptoListingsLatestResponse,
) (map[string]float64, error) {
	cryptocurrency := s.GetCryptocurrencyByID(crypts, int(portfolioDetail.CMCCryptocurrencyID))
	if cryptocurrency == nil {
		return nil, ErrNotFound
	}

	sum := portfolioDetail.Price * portfolioDetail.Count
	nowSum := cryptocurrency.Quote.USD.Price * portfolioDetail.Count
	additionalData := make(map[string]float64)
	additionalData["sum"] = sum
	additionalData["now_price"] = cryptocurrency.Quote.USD.Price
	additionalData["now_sum"] = nowSum
	additionalData["percent_change_24h"] = cryptocurrency.Quote.USD.PercentChange24h
	additionalData["percent_change_30d"] = cryptocurrency.Quote.USD.PercentChange30d
	additionalData["percent_change_90d"] = cryptocurrency.Quote.USD.PercentChange90d
	additionalData["portfolio_share"] = (cryptocurrency.Quote.USD.Price * portfolioDetail.Count) / totalNowSum
	additionalData["ROI"] = ((nowSum - sum) / sum) * 100
	additionalData["profit"] = nowSum - sum
	return additionalData, nil
}

func (s *Service) GetCryptocurrencyByID(crypts *clients.CryptoListingsLatestResponse, cryptID int) *clients.CryptocurrencyDataDetail {
	data := crypts.Data
	for len(data) > 0 {
		curr := len(data) / 2
		currCrypt := data[curr]
		currCryptID := currCrypt.ID

		if currCryptID == cryptID {
			return &currCrypt
		} else if currCryptID > cryptID {
			data = data[:curr]
		} else {
			data = data[curr+1:]
		}
	}

	return nil
}
