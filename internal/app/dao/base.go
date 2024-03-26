package dao

import (
	core "fl/my-portfolio/internal/app"
	"fl/my-portfolio/internal/app/models"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

type DAO struct {
	*sqlx.DB
    logger *logrus.Logger
}

func NewDAO() *DAO {
	return &DAO{core.GetDB(), core.GetLogger()}
}

func (d *DAO) GetOne(model models.IModel, objectID uuid.UUID) error {
    query, args, err := psql.Select("*").From(model.TableName()).Where(sq.Eq{"id": objectID}).ToSql()
	if err != nil {
        return err
	}

	if err := d.Get(model, query, args...); err != nil {
        return err
	}

    return nil
}

func (d *DAO) GetAll(model models.IModel) ([]models.IModel, error) {
    var objects []models.IModel
    query, _, err := psql.Select("*").From(model.TableName()).ToSql()
    if err != nil {
        return nil, err
    }

    if err := d.Select(objects, query); err != nil {
        return nil , err

    }

    return objects, nil
}


func (d *DAO) Delete(model models.IModel, objectId uuid.UUID) error {
    query, args, err := psql.Delete(model.TableName()).Where(sq.Eq{"id": objectId}).ToSql()
    if err != nil {
        return err
    }

    if _, err := d.Exec(query, args...); err != nil {
        return err
    }

    return nil
}
