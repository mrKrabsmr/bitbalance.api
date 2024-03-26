package controllers

import "net/http"

// Get @Summary
// @Description все криптовалюты
// @Tags cryptocurrency 
// @Accept json
// @Produce json
// @Success 200 {object} ResponseSuccess
// @Failure 500 {object} ResponseError
// @Security ApiKeyAuth
// @Router /api/v1/cryptocurrencies [get]
func (c *Controller) Cryptocurrencies(writer http.ResponseWriter, request *http.Request) {
    cryptocurrencies, err := c.service.GetAllCryptocurrencies(false)
    if err != nil {
        c.logger.Error(err)
        c.JSONResponse(writer, text[0], http.StatusInternalServerError)
        return
    }

    c.JSONResponse(writer, cryptocurrencies, http.StatusOK)
}
