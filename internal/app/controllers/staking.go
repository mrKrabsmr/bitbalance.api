package controllers

import (
	"fl/my-portfolio/internal/app/models"
	"net/http"

	"github.com/go-chi/chi"
)

// Get @Summary
// @Description список стейкингов
// @Tags staking 
// @Accept json
// @Produce json
// @Success 200 {object} []StakingGetResponse
// @Failure 500 {object} ResponseError
// @Security ApiKeyAuth
// @Router /api/v1/stakings [get]
func (c *Controller) Stakings(writer http.ResponseWriter, request *http.Request) {
	crypts, err := c.service.GetAllCryptocurrencies(false)
	if err != nil {
		c.logger.Error(err)
		c.JSONResponse(writer, text[0], http.StatusInternalServerError)
		return
	}

	binanceStakings, err := c.service.GetAllStakingData("binance")
	// bybitStakings, err := c.service.GetAllStakingData("bybit")
	okxStakings, err := c.service.GetAllStakingData("okx")
	if err != nil {
		c.logger.Error(err)
		c.JSONResponse(writer, text[0], http.StatusInternalServerError)
		return
	}

	cryptStakingList := make([]StakingGetResponse, 0, len(crypts.Data))
	for _, crypt := range crypts.Data {
		binanceStaking := c.service.GetStakingDataByCryptSymbol(binanceStakings, crypt.Symbol)
		// bybitStaking := c.service.GetStakingDataByCryptSymbol(bybitStakings, crypt.Symbol)
		okxStaking := c.service.GetStakingDataByCryptSymbol(okxStakings, crypt.Symbol)

		cryptStaking := StakingGetResponse{
			CrpyptocurrencyData: crypt,
			BinanceStaking:      binanceStaking,
			// BybitStaking:        bybitStaking,
			OkxStaking:          okxStaking,
		}

		cryptStakingList = append(cryptStakingList, cryptStaking)
	}

	c.JSONResponse(writer, cryptStakingList, http.StatusOK)
}

// Get @Summary
// @Description список стейкингов для криптовалют в портфеле
// @Tags staking 
// @Accept json
// @Produce json
// @Success 200 {object} []StakingGetResponse
// @Failure 500 {object} ResponseError
// @Security ApiKeyAuth
// @Router /api/v1/stakings/portfolio [get]
func (c *Controller) StakingsPortfolio(writer http.ResponseWriter, request *http.Request) {
	crypts, err := c.service.GetAllCryptocurrencies(true)
	if err != nil {
		c.logger.Error(err)
		c.JSONResponse(writer, text[0], http.StatusInternalServerError)
		return
	}

	binanceStakings, err := c.service.GetAllStakingData("binance")
	// bybitStakings, err := c.service.GetAllStakingData("bybit")
	okxStakings, err := c.service.GetAllStakingData("okx")
	if err != nil {
		c.logger.Error(err)
		c.JSONResponse(writer, text[0], http.StatusInternalServerError)
		return
	}

	user := request.Context().Value("user").(*models.User)

	portfolio, err := c.service.GetPortfolioCryptocurrencies(user.ID)
	if err != nil {
		c.logger.Error(err)
		c.JSONResponse(writer, text[0], http.StatusInternalServerError)
		return
	}

	cryptStakingList := make([]StakingGetResponse, 0, len(crypts.Data))
	for _, portfolioDetail := range portfolio {
		crypt := c.service.GetCryptocurrencyByID(crypts, int(portfolioDetail.CMCCryptocurrencyID))
		binanceStaking := c.service.GetStakingDataByCryptSymbol(binanceStakings, crypt.Symbol)
		// bybitStaking := c.service.GetStakingDataByCryptSymbol(bybitStakings, crypt.Symbol)
		okxStaking := c.service.GetStakingDataByCryptSymbol(okxStakings, crypt.Symbol)

		cryptStaking := StakingGetResponse{
			CrpyptocurrencyData: *crypt,
			BinanceStaking:      binanceStaking,
			// BybitStaking:        bybitStaking,
			OkxStaking:          okxStaking,
		}

		cryptStakingList = append(cryptStakingList, cryptStaking)
	}

	c.JSONResponse(writer, cryptStakingList, http.StatusOK)
}

// Get @Summary
// @Description список стейкингов для конкретной криптовалюты
// @Tags staking 
// @Accept json
// @Produce json
// @Param crypt_symbol  path string true "Cryptocurrency symbol"
// @Success 200 {object} []StakingDetailGetResponse
// @Failure 500 {object} ResponseError
// @Security ApiKeyAuth
// @Router /api/v1/stakings/detail/{crypt_symbol} [get]
func (c *Controller) StakingsDetail(writer http.ResponseWriter, request *http.Request) {
	cryptSymbol := chi.URLParam(request, "crypt_symbol")
	binanceStakings, err := c.service.GetAllStakingData("binance")
	// bybitStakings, err := c.service.GetAllStakingData("bybit")
	okxStakings, err := c.service.GetAllStakingData("okx")
	if err != nil {
		c.logger.Error(err)
		c.JSONResponse(writer, text[0], http.StatusInternalServerError)
		return
	}

	binanceStakingsFound := c.service.GetAllStakingDataByCryptSymbol(binanceStakings, cryptSymbol)
	// bybitStakingsFound := c.service.GetAllStakingDataByCryptSymbol(bybitStakings, cryptSymbol)
	okxStakingsFound := c.service.GetAllStakingDataByCryptSymbol(okxStakings, cryptSymbol)
	stakingsFound := &StakingDetailGetResponse{
		BinanceStaking: binanceStakingsFound,
		// BybitStaking:   bybitStakingsFound,
		OkxStaking:     okxStakingsFound,
	}

	c.JSONResponse(writer, stakingsFound, http.StatusOK)
}
