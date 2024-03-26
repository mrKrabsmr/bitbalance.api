package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	GenderMail   = "mail"
	GenderFemail = "femail"
)

type User struct {
	ID          uuid.UUID  `db:"id" json:"id"`
	Email       string     `db:"email" json:"email"`
	Password    string     `db:"password" json:"password"`
	FirstName   string     `db:"first_name" json:"first_name"`
	LastName    string     `db:"last_name" json:"last_name"`
	Gender      string     `db:"gender" json:"gender"`
	BirthDate   time.Time `db:"birth_date" json:"birth_date"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	IsSuperUser bool       `db:"is_superuser" json:"is_superuser"`
	IsStaff     bool       `db:"is_staff" json:"is_staff"`
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
