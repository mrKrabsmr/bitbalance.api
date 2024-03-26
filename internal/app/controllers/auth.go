package controllers

import (
	"encoding/json"
	"errors"
	core "fl/my-portfolio/internal/app"
	"fl/my-portfolio/internal/app/models"
	"fl/my-portfolio/internal/app/services"
	"io"
	"net/http"
	"time"
)

// Post @Summary
// @Description регистрация
// @Tags auth 
// @Accept json
// @Produce json
// @Param input body RegisterData true "data"
// @Success 201 {object} Tokens
// @Failure 400,500 {object} ResponseError
// @Router /api/v1/register [post]
func (c *Controller) Register(writer http.ResponseWriter, request *http.Request) {
	var obj *RegisterData

	defer request.Body.Close()
	data, err := io.ReadAll(request.Body)
	if err != nil {
		c.logger.Error(err)
		c.JSONResponse(writer, text[0], http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(data, &obj); err != nil {
		c.logger.Error(err)
		c.JSONResponse(writer, err.Error(), http.StatusBadRequest)
		return
	}

	if err := core.GetValidator().Struct(obj); err != nil {
		c.logger.Error(err)
		c.JSONResponse(writer, err.Error(), http.StatusBadRequest)
		return
	}

	user := &models.User{
		Email:     obj.Email,
		Password:  obj.Password,
		FirstName: obj.FirstName,
		LastName:  obj.LastName,
		Gender:    obj.Gender,
		BirthDate: time.Time(obj.BirthDate),
	}

	userID, err := c.service.CreateUser(user)
	if err != nil {
		var errTxt string
		var errCode int
		if errors.Is(err, services.ErrDuplicate) {
			errTxt = "user with such email already exists"
			errCode = http.StatusBadRequest
		} else {
			errTxt = text[0]
			errCode = http.StatusInternalServerError
		}

		c.logger.Error(err)
		c.JSONResponse(writer, errTxt, errCode)
		return
	}

	accessToken, refreshToken, err := c.service.GenerateTokens(userID)
	if err != nil {
		c.logger.Error(err)
		c.JSONResponse(writer, text[0], http.StatusInternalServerError)
		return
	}

	c.JSONResponse(
		writer,
		Tokens{
			UserID:       userID.String(),
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
		http.StatusCreated,
	)
}

// Post @Summary
// @Description авторизация
// @Tags auth 
// @Accept json
// @Produce json
// @Param input body LoginData true "data"
// @Success 200 {object} Tokens
// @Failure 400,500 {object} ResponseError
// @Router /api/v1/login [post]
func (c *Controller) Login(writer http.ResponseWriter, request *http.Request) {
	var obj *LoginData

	defer request.Body.Close()
	data, err := io.ReadAll(request.Body)
	if err != nil {
		c.logger.Error(err)
		c.JSONResponse(writer, text[0], http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(data, &obj); err != nil {
		c.logger.Error(err)
		c.JSONResponse(writer, err.Error(), http.StatusBadRequest)
		return
	}

	if err := core.GetValidator().Struct(obj); err != nil {
		c.logger.Error(err)
		c.JSONResponse(writer, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := c.service.GetUserLogin(obj.Email, obj.Password)
	if err != nil {
		c.logger.Error(err)
		c.JSONResponse(writer, "incorrect email - password pair", http.StatusBadRequest)
		return
	}

	accessToken, refreshToken, err := c.service.GenerateTokens(user.ID)
	if err != nil {
		c.logger.Error(err)
		c.JSONResponse(writer, text[0], http.StatusInternalServerError)
		return
	}

	c.JSONResponse(
		writer,
		Tokens{
			UserID:       user.ID.String(),
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
		http.StatusOK,
	)
}

// Post @Summary
// @Description обновление токенов
// @Tags auth 
// @Accept json
// @Produce json
// @Param input body RefreshData true "data"
// @Success 200 {object} Tokens
// @Failure 400,500 {object} ResponseError
// @Router /api/v1/refresh [post]
func (c *Controller) Refresh(writer http.ResponseWriter, request *http.Request) {
	var obj *RefreshData

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

	_, err = c.service.ParseJWT(obj.RefreshToken)
	if err != nil {
		c.logger.Error(err)
		c.JSONResponse(writer, "invalid refresh token", http.StatusBadRequest)
		return
	}

	user, err := c.service.GetUserByRefreshToken(obj.RefreshToken)
	if err != nil {
		c.logger.Error(err)
		c.JSONResponse(writer, "user not found", http.StatusBadRequest)
		return
	}

	accessToken, refreshToken, err := c.service.GenerateTokens(user.ID)
	if err != nil {
		c.logger.Error(err)
		c.JSONResponse(writer, text[0], http.StatusInternalServerError)
		return
	}

	c.JSONResponse(
		writer,
		Tokens{
			UserID:       user.ID.String(),
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
		http.StatusOK,
	)
}
