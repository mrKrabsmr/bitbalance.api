package controllers

import (
	"encoding/json"
	"errors"
	core "fl/my-portfolio/internal/app"
	"fl/my-portfolio/internal/app/models"
	"fl/my-portfolio/internal/app/services"
	"fl/my-portfolio/pkg/types"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

// Get @Summary
// @Description портфель
// @Tags portfolio
// @Accept json
// @Produce json
// @Success 200 {object} PortfolioGetResponse
// @Failure 500 {object} ResponseError
// @Security ApiKeyAuth
// @Router /api/v1/portfolio [get]
func (c *Controller) GetPortfolio(writer http.ResponseWriter, request *http.Request) {
	user := request.Context().Value("user").(*models.User)

	portfolio, err := c.service.GetPortfolioCryptocurrencies(user.ID)
	if err != nil {
		c.logger.Error(err)
		c.JSONResponse(writer, text[0], http.StatusInternalServerError)
		return
	}

	cryptocurrencies, err := c.service.GetAllCryptocurrencies(true)
	if err != nil {
		c.logger.Error(err)
		c.JSONResponse(writer, text[0], http.StatusInternalServerError)
		return
	}

	totalPurchaseSum, totalNowSum, err := c.service.GetSumPortfolioCryptocurrencies(portfolio)
	if err != nil {
		c.logger.Error(err)
		c.JSONResponse(writer, text[0], http.StatusInternalServerError)
		return
	}

	cryptocurrenciesData := make([]PortfolioDetailGetResponse, 0, len(portfolio))
	for _, portfolioDetail := range portfolio {
		additionalData, err := c.service.GetCryptocurrencyAdditionalData(portfolioDetail, totalNowSum, cryptocurrencies)
		if err != nil {
			c.logger.Error(err)
			c.JSONResponse(writer, text[0], http.StatusInternalServerError)
			return
		}

		detailData := PortfolioDetailGetResponse{
			ID:                   portfolioDetail.ID,
			CryptID:              portfolioDetail.CMCCryptocurrencyID,
			Cryptocurrency:       portfolioDetail.Cryptocurrency,
			CryptocurrencySymbol: portfolioDetail.CryptocurrencySymbol,
			Price:                portfolioDetail.Price,
			Count:                portfolioDetail.Count,
			Sum:                  additionalData["sum"],
			NowPrice:             additionalData["now_price"],
			NowSum:               additionalData["now_sum"],
			PurchaseTime:         types.Time(portfolioDetail.PurchaseTime),
			Commentary:           portfolioDetail.Commentary,
			CreatedAt:            types.Time(portfolioDetail.CreatedAt),
			PercentChange24h:     additionalData["percent_change_24h"],
			PercentChange30d:     additionalData["percent_change_30d"],
			PercentChange90d:     additionalData["percent_change_90d"],
			PortfolioShare:       additionalData["portfolio_share"],
			ROI:                  additionalData["ROI"],
		}

		cryptocurrenciesData = append(cryptocurrenciesData, detailData)
	}

	portfolioData := &PortfolioGetResponse{
		TotalPurchaseSum: totalPurchaseSum,
		TotalNowSum:      totalNowSum,
		ROI:              ((totalNowSum - totalPurchaseSum) / totalPurchaseSum) * 100,
		Cryptocurrencies: cryptocurrenciesData,
	}

	c.JSONResponse(writer, portfolioData, http.StatusOK)
}

// Post @Summary
// @Description добавить в портфель криптовалюту
// @Tags portfolio
// @Accept json
// @Produce json
// @Param input body PortfolioPostData true "data"
// @Success 200 {object} ResponseSuccess 
// @Failure 500 {object} ResponseError
// @Security ApiKeyAuth
// @Router /api/v1/portfolio [post]
func (c *Controller) CreatePortfolio(writer http.ResponseWriter, request *http.Request) {
	var obj *PortfolioPostData

	defer request.Body.Close()
	data, err := io.ReadAll(request.Body)
	if err != nil {
		c.logger.Error(err)
		c.JSONResponse(writer, text[0], http.StatusInternalServerError)
		return
	}

	if err = json.Unmarshal(data, &obj); err != nil {
		c.logger.Error(err)
		c.JSONResponse(writer, err.Error(), http.StatusBadRequest)
		return
	}

	if err = core.GetValidator().Struct(obj); err != nil {
		c.logger.Error(err)
		c.JSONResponse(writer, err.Error(), http.StatusBadRequest)
		return
	}

	user := request.Context().Value("user").(*models.User)

	portfolio := &models.Portfolio{
		UserID:              user.ID,
		CMCCryptocurrencyID: obj.CMCCryptocurrencyID,
		Price:               obj.Price,
		Count:               obj.Count,
		PurchaseTime:        time.Time(obj.PurchaseTime),
		Commentary:          obj.Commentary,
	}

	if err = c.service.AddCryptocurrencyToPortfolio(portfolio); err != nil {
		c.logger.Error(err)
		if errors.Is(err, services.ErrNotFound) {
			c.JSONResponse(writer, "cryptocurrency with such id was not found", http.StatusNotFound)
			return
		}

		if errors.Is(err, services.ErrDuplicate) {
			c.JSONResponse(writer, "cryptocurrency has already been added to the portfolio", http.StatusBadRequest)
			return
		}
		c.JSONResponse(writer, text[0], http.StatusInternalServerError)
		return
	}

	c.JSONResponse(writer, "success added", http.StatusCreated)
}

// Put,Patch @Summary
// @Description обновитель информацию о криптовалюте в портфеле
// @Tags portfolio
// @Accept json
// @Produce json
// @Param input body PortfolioPatchData true "data"
// @Param id  path string true "Portfolio detail ID"
// @Success 200 {object} ResponseSuccess
// @Failure 500 {object} ResponseError
// @Security ApiKeyAuth
// @Router /api/v1/portfolio/{id} [patch]
func (c *Controller) UpdatePortfolioDetail(writer http.ResponseWriter, request *http.Request) {
	idStr := chi.URLParam(request, "id")
	portfolioDetailID, err := uuid.Parse(idStr)
	if err != nil {
		c.logger.Error(err)
		c.JSONResponse(writer, err.Error(), http.StatusBadRequest)
		return
	}

	var obj *PortfolioPatchData

	defer request.Body.Close()
	data, err := io.ReadAll(request.Body)
	if err != nil {
		c.logger.Error(err)
		c.JSONResponse(writer, text[0], http.StatusInternalServerError)
		return
	}

	if err = json.Unmarshal(data, &obj); err != nil {
		c.logger.Error(err)
		c.JSONResponse(writer, err.Error(), http.StatusBadRequest)
		return
	}

	if err = core.GetValidator().Struct(obj); err != nil {
		c.logger.Error(err)
		c.JSONResponse(writer, err.Error(), http.StatusBadRequest)
		return
	}

	portfolio := &models.Portfolio{
		Price:        obj.Price,
		Count:        obj.Count,
		PurchaseTime: time.Time(obj.PurchaseTime),
		Commentary:   obj.Commentary,
	}

	if err = c.service.UpdateCryptocurrencyDataInPortfolio(portfolioDetailID, portfolio, request.Method == "PATCH"); err != nil {
		c.logger.Error(err)
		c.JSONResponse(writer, text[0], http.StatusInternalServerError)
		return
	}

	c.JSONResponse(writer, "success updated", http.StatusCreated)
}

// Delete @Summary
// @Description удалить криптовалюту из портфеля
// @Tags portfolio
// @Accept json
// @Produce json
// @Success 200 {object} PortfolioGetResponse
// @Failure 500 {object} ResponseError
// @Param id  path string true "Portfolio detail ID"
// @Security ApiKeyAuth
// @Router /api/v1/portfolio/{id} [delete]
func (c *Controller) DeletePortfolioDetail(writer http.ResponseWriter, request *http.Request) {
	idStr := chi.URLParam(request, "id")
	portfolioDetailID, err := uuid.Parse(idStr)
	if err != nil {
		c.logger.Error(err)
		c.JSONResponse(writer, err.Error(), http.StatusBadRequest)
		return
	}

	if err = c.service.DeleteCryptocurrencyFromPortfolio(portfolioDetailID); err != nil {
		c.logger.Error(err)
		c.JSONResponse(writer, err.Error(), http.StatusBadRequest)
		return
	}

	c.JSONResponse(writer, "success deleted", http.StatusCreated)
}
