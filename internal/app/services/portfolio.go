package services

import (
	core "fl/my-portfolio/internal/app"
	"fl/my-portfolio/internal/app/models"
	"time"

	"github.com/google/uuid"
)

func (s *Service) AddCryptocurrencyToPortfolio(portfolio *models.Portfolio) error {
	crypts, err := s.GetAllCryptocurrencies(true)
	if err != nil {
		return nil
	}

	cryptocurrency := s.GetCryptocurrencyByID(crypts, int(portfolio.CMCCryptocurrencyID))
	if cryptocurrency == nil {
		return ErrNotFound
	}

	pID, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	portfolio.ID = pID
	portfolio.CreatedAt = time.Now().In(core.GetLocation())
	portfolio.Cryptocurrency = cryptocurrency.Name
	portfolio.CryptocurrencySymbol = cryptocurrency.Symbol

	if err = s.dao.CreatePortfolioDetail(portfolio); err != nil {
		return getError(err)
	}

	return nil
}

func (s *Service) UpdateCryptocurrencyDataInPortfolio(
	portfolioDetailID uuid.UUID, portfolio *models.Portfolio, patch bool) error {
	obj := &models.Portfolio{}
	if err := s.dao.GetOne(obj, portfolioDetailID); err != nil {
		return getError(err)
	}

	if patch {
		if portfolio.Commentary == "" {
			portfolio.Commentary = obj.Commentary
		}
		if portfolio.Count == 0 {
			portfolio.Count = obj.Count
		}
		if portfolio.Price == 0 {
			portfolio.Price = obj.Price
		}
		if portfolio.PurchaseTime.IsZero() {
			portfolio.PurchaseTime = obj.PurchaseTime
		}
	}

	if err := s.dao.UpdatePortfolioDetail(portfolioDetailID, portfolio); err != nil {
		return getError(err)
	}

	return nil
}

func (s *Service) DeleteCryptocurrencyFromPortfolio(portfolioDetailId uuid.UUID) error {
	if err := s.dao.DeletePortfolioDetail(portfolioDetailId); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetPortfolioCryptocurrencies(userID uuid.UUID) ([]*models.Portfolio, error) {
	portfolio, err := s.dao.GetUserPortfolio(userID)
	if err != nil {
		return nil, getError(err)
	}

	return portfolio, nil
}

func (s *Service) GetSumPortfolioCryptocurrencies(portfolio []*models.Portfolio) (float64, float64, error) {
	var totalPurchaseSum, totalNowSum float64

	crypts, err := s.GetAllCryptocurrencies(true)
	if err != nil {
		return -1, -1, err
	}

	for _, portfolioDetail := range portfolio {
		totalPurchaseSum += portfolioDetail.Price * portfolioDetail.Count

		cryptocurrency := s.GetCryptocurrencyByID(crypts, int(portfolioDetail.CMCCryptocurrencyID))
		if cryptocurrency == nil {
			return -1, -1, ErrNotFound
		}

		totalNowSum += cryptocurrency.Quote.USD.Price * portfolioDetail.Count
	}

	return totalPurchaseSum, totalNowSum, nil
}
