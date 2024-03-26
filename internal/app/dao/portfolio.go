package dao

import (
	"fl/my-portfolio/internal/app/models"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

func (d *DAO) CreatePortfolioDetail(portfolio *models.Portfolio) error {
	query, args, err := psql.Insert(portfolio.TableName()).Columns(
		"id", "user_id", "cmc_cryptocurrency_id", "cryptocurrency", "cryptocurrency_symbol", "price", "count",
		"purchase_time", "commentary", "created_at",
	).Values(
		portfolio.ID, portfolio.UserID, portfolio.CMCCryptocurrencyID, portfolio.Cryptocurrency, portfolio.CryptocurrencySymbol,
		portfolio.Price, portfolio.Count, portfolio.PurchaseTime, portfolio.Commentary, portfolio.CreatedAt,
	).ToSql()
	if err != nil {
		return err
	}

	if _, err = d.Exec(query, args...); err != nil {
		return err
	}

	return nil
}

func (d *DAO) GetUserPortfolio(userID uuid.UUID) ([]*models.Portfolio, error) {
	var portfolio []*models.Portfolio

	query, args, err := psql.Select("*").From("portfolios").Where(sq.Eq{"user_id": userID}).ToSql()
	if err != nil {
		return nil, err
	}

	if err = d.Select(&portfolio, query, args...); err != nil {
		return nil, err
	}

	return portfolio, nil
}

func (d *DAO) UpdatePortfolioDetail(id uuid.UUID, portfolio *models.Portfolio) error {
	query, args, err := psql.Update(portfolio.TableName()).Set("price", portfolio.Price).Set("count", portfolio.Count,
	).Set("purchase_time", portfolio.PurchaseTime).Set("commentary", portfolio.Commentary).Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return err
	}

	if _, err = d.Exec(query, args...); err != nil {
		return err
	}

	return nil
}

func (d *DAO) DeletePortfolioDetail(id uuid.UUID) error {
	query, args, err := psql.Delete("portfolios").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return err
	}

	if _, err = d.Exec(query, args...); err != nil {
		return err
	}

	return nil
}
