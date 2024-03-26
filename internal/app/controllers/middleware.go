package controllers

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func (c *Controller) AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		tokenBearer := request.Header.Get("Authorization")
		s := strings.Split(tokenBearer, " ")
		if len(s) != 2 {
			c.JSONResponse(writer, "invalid or empty authentication data", http.StatusUnauthorized)
			return
		}

		if s[0] != "Bearer" {
			c.JSONResponse(writer, "authentication data must be in format 'Bearer `token`'", http.StatusUnauthorized)
			return
		}

		t := s[1]
		token, err := c.service.ParseJWT(t)
		if err != nil {
			c.logger.Error(err)
			c.JSONResponse(writer, "token is invalid or expired", http.StatusUnauthorized)
			return
		}

		claims, _ := token.Claims.(jwt.MapClaims)
		userIDStr, ok := claims["user_id"].(string)

		user, err := c.service.GetUser(userIDStr)
		if !ok || err != nil {
			c.logger.Error(err)
			c.JSONResponse(writer, "invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(request.Context(), "user", user)

		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}
