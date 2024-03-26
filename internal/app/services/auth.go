package services

import (
	"errors"
	core "fl/my-portfolio/internal/app"
	"fl/my-portfolio/internal/app/models"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	accessTokenMaxAge  = time.Minute * 15 
	refreshTokenMaxAge = time.Hour * 24 * 30
)

func (s *Service) CreateUser(user *models.User) (uuid.UUID, error) {
	uID, err := uuid.NewRandom()
	if err != nil {
		s.logger.Error(err)
		return uuid.UUID{}, err
	}

	user.ID = uID
	user.CreatedAt = time.Now()

	password, err := s.hashPassword(user.Password)
	if err != nil {
		s.logger.Error(err)
		return uuid.UUID{}, err
	}

	user.Password = "bcrypt$" + string(password)

	err = s.dao.CreateUser(user)
	if err != nil {
		s.logger.Error(err)
		return uuid.UUID{}, getError(err)
	}

	return uID, nil
}

func (s *Service) GetUserLogin(email string, password string) (*models.User, error) {
	user, err := s.dao.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	activePassword, _ := strings.CutPrefix(user.Password, "bcrypt$")
	if err = bcrypt.CompareHashAndPassword([]byte(activePassword), []byte(password)); err != nil {
		return nil, err
	}

	return user, err
}

func (s *Service) GetUserByRefreshToken(refreshToken string) (*models.User, error) {
	return s.dao.GetUserByRefreshToken(refreshToken)
}

func (s *Service) hashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), 8)
}

func (s *Service) GenerateTokens(userID uuid.UUID) (string, string, error) {
	timeNow := time.Now().In(core.GetLocation())
	firstToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"token_type": "access",
		"exp":        jwt.NewNumericDate(timeNow.Add(accessTokenMaxAge)),
		"iss":        jwt.NewNumericDate(timeNow),
		"user_id":    userID.String(),
	})

	accessToken, err := firstToken.SignedString(s.key)
	if err != nil {
		return "", "", err
	}

	secondToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"token_type": "refresh",
		"exp":        jwt.NewNumericDate(timeNow.Add(refreshTokenMaxAge)),
		"iss":        jwt.NewNumericDate(timeNow),
	})

	refreshToken, err := secondToken.SignedString(s.key)
	if err != nil {
		return "", "", err
	}

	sID, err := uuid.NewRandom()
	if err != nil {
		s.logger.Error(err)
		return "", "", err
	}

	s.dao.CreateSession(&models.Session{
		ID:           sID,
		UserID:       userID,
		RefreshToken: refreshToken,
	})

	return accessToken, refreshToken, nil
}

func (s *Service) ParseJWT(tokenStr string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenStr, func(tokenStr *jwt.Token) (interface{}, error) {
		if _, ok := tokenStr.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return s.key, nil
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}