package controllers

import (
	"encoding/json"
	core "fl/my-portfolio/internal/app"
	"fl/my-portfolio/internal/app/services"
	"net/http"

	"github.com/sirupsen/logrus"
)

type Controller struct {
	service *services.Service
	logger  *logrus.Logger
}

func NewController() *Controller {
	return &Controller{
		service: services.NewService(),
		logger:  core.GetLogger(),
	}
}

func (c *Controller) JSONResponse(writer http.ResponseWriter, data interface{}, statusCode int) {
	var jsonData []byte

	if statusCode >= 400 {
		resp := new(ResponseError)
		resp.Error = true
		resp.Message = data.(string)

		jsonData, _ = json.Marshal(resp)
	} else {
		resp := new(ResponseSuccess)
		resp.Error = false
		resp.Data = data

		jD, err := json.Marshal(resp)
		if err != nil {
			c.logger.Info(err)
			return
		}

		jsonData = jD
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	writer.Write(jsonData)
}

func (c *Controller) JSONResponseBinary(writer http.ResponseWriter, jsonData []byte, statusCode int) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	writer.Write(jsonData)
}
