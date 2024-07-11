package controllers

import (
	"fl/my-portfolio/internal/app/models"
	"net/http"
)

func (c *Controller) GetMe(writer http.ResponseWriter, request *http.Request) {
	user := request.Context().Value("user").(*models.User)
	user.Password = ""
    
    c.JSONResponse(writer, user, http.StatusOK)
}
