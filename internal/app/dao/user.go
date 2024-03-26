package dao

import (
	"fl/my-portfolio/internal/app/models"

	sq "github.com/Masterminds/squirrel"
)


func (d *DAO) CreateUser(user *models.User) error {
    query, args, err := psql.Insert(user.TableName()).Columns(
        "id", "email", "password", "first_name", "last_name", "gender", "birth_date", "created_at",
    ).Values(
        user.ID, user.Email, user.Password, user.FirstName, user.LastName, user.Gender, user.BirthDate, user.CreatedAt,
    ).ToSql()

    if err != nil {
        return err
    }

    if _, err := d.Exec(query, args...); err != nil {
        return err
    }

    return nil
}

func (d *DAO) CreateSession(session *models.Session) error {
    query, args, err := psql.Insert(session.TableName()).Columns(
        "id", "user_id", "refresh_token", 
    ).Values(
        session.ID, session.UserID, session.RefreshToken,
    ).ToSql()
    if err != nil {
        return err
    }

    if _, err := d.Exec(query, args...); err != nil {
        return err
    }

    return nil
}

func (d *DAO) GetUserByEmail(email string) (*models.User, error) {
    var user models.User

    query, args, err := psql.Select("*").From("users").Where(sq.Eq{"email": email}).ToSql()
    if err != nil {
        return nil, err
    }

    if err := d.Get(&user, query, args...); err != nil {
        return nil, err
    }

    return &user, nil
}

func (d *DAO) GetUserByRefreshToken(refreshToken string) (*models.User, error) {
    var user models.User

    query, args, err := psql.Select("users.*").From(user.TableName(),
        ).Join("sessions ON (users.id = sessions.user_id)").Where(sq.Eq{"refresh_token": refreshToken}).ToSql()
    if err != nil {
        return nil, err
    }

    if err = d.Get(&user, query, args...); err != nil {
        return nil, err
    }

    return &user, nil
}
