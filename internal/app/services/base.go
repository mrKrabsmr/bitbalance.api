package services

import (
	"errors"
	core "fl/my-portfolio/internal/app"
	"fl/my-portfolio/internal/app/cache"
	"fl/my-portfolio/internal/app/dao"
	"fl/my-portfolio/internal/clients"
	"fl/my-portfolio/internal/configs"

	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type Service struct {
	dao    *dao.DAO
	logger *logrus.Logger
	config *configs.Config
	key    []byte

	cmcClient      *clients.CMCClient
	exchangeClient *clients.ExchangeClient
	cacher         *cache.Redis
}

func NewService() *Service {
	return &Service{
		dao:            dao.NewDAO(),
		logger:         core.GetLogger(),
		config:         core.GetConfig(),
		key:            core.GetKey(),
		cmcClient:      clients.NewCMCClient(),
		exchangeClient: clients.NewExchangeClient(),
		cacher:         cache.NewRedis(),
	}
}

func getError(err error) error {
	var pgErr *pq.Error
	ok := errors.As(err, &pgErr)
	if ok && pgErr.Code == "23505" {
		return ErrDuplicate
	}

	return err

}
