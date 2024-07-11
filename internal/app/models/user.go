package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID `db:"id" json:"id"`
	Username    string    `db:"username" json:"username"`
	Password    string    `db:"password" json:"password"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	IsSuperUser bool      `db:"is_superuser" json:"is_superuser"`
	IsStaff     bool      `db:"is_staff" json:"is_staff"`
}

func (u *User) TableName() string {
	return "users"
}

type Session struct {
	ID           uuid.UUID `db:"id" json:"id"`
	UserID       uuid.UUID `db:"user_id" json:"user_id"`
	RefreshToken string    `db:"refresh_token" json:"refresh_token"`
}

func (s *Session) TableName() string {
	return "sessions"
}
