package services

import (
	"context"
	"encoding/json"
	"errors"
	"fl/my-portfolio/internal/clients"
	"sort"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

func (s *Service) GetAllStakingData(exch string) ([]clients.StakingData, error) {
	stakings, err := s.cacher.Get(context.Background(), exch + "_stakings").Result()
	if err != nil {
		s.logger.Info(err)
		if errors.Is(err, redis.Nil) {
			binanceSt, err := s.exchangeClient.GetBinanceStakingData()
			if err != nil {
				return nil, err
			}
			sort.Slice(binanceSt, func(i, j int) bool {
				return binanceSt[i].Asset < binanceSt[j].Asset
			})
			binanceJsonData, err := json.Marshal(binanceSt)
			if err != nil {
				return nil, err
			}

			okxSt, err := s.exchangeClient.GetOkxStakingData()
			if err != nil {
				return nil, err
			}
			sort.Slice(okxSt, func(i, j int) bool {
				return okxSt[i].Asset < okxSt[j].Asset
			})
			okxJsonData, err := json.Marshal(okxSt)
			if err != nil {
				return nil, err
			}

			go func() {
				s.cacher.Set(context.Background(), "binance_stakings", binanceJsonData, time.Minute*60)
				s.cacher.Set(context.Background(), "okx_stakings", okxJsonData, time.Minute*60)
			}()
			
			switch exch {
			case "binance":
				return binanceSt, nil
			case "okx":
				return okxSt, nil
			default:
				return nil, nil
			}
		} else {
			return nil, err
		}
	}

	var objects []clients.StakingData

	jsonString, err := json.Marshal(stakings)
	if err != nil {
		return nil, err
	}

	ss, _ := strconv.Unquote(string(jsonString))
	if err = json.Unmarshal([]byte(ss), &objects); err != nil {
		return nil, err
	}

	return objects, nil
}

func (s *Service) GetStakingDataByCryptSymbol(stakingData []clients.StakingData, cryptSymbol string) *clients.StakingData {
	data := stakingData
	for len(data) > 0 {
		curr := len(data) / 2
		currEl := data[curr]
		if currEl.Asset == cryptSymbol {
			return &currEl
		} else if currEl.Asset > cryptSymbol {
			data = data[:curr]
		} else {
			data = data[curr+1:]
		}
	}

	return nil 
}

func (s *Service) GetAllStakingDataByCryptSymbol(stakingData []clients.StakingData, cryptSymbol string) []clients.StakingData {
	data := stakingData
	for len(data) > 0 {
		curr := len(data) / 2
		currEl := data[curr]
		if currEl.Asset == cryptSymbol {
			max, min := curr, curr
			for ; ; max++ {
				if max == len(data) || data[max].Asset != cryptSymbol {
					break
				}
			}

			for ; ; min-- {
				if min < 0 || data[min].Asset != cryptSymbol {
					break
				}
			}

			return data[min+1 : max]

		} else if currEl.Asset > cryptSymbol {
			data = data[:curr]
		} else {
			data = data[curr+1:]
		}
	}

	return nil
}
