package services

import (
	"fl/my-portfolio/internal/app/models"

	"github.com/google/uuid"
)


func (s *Service) GetUser(userIDStr string) (*models.User, error) {
    userID, err := uuid.Parse(userIDStr)
    if err != nil {
        return nil, err
    }
    
    user := &models.User{}
    if err = s.dao.GetOne(user, userID); err != nil {
        return nil, err
    }

    return user, nil
}
