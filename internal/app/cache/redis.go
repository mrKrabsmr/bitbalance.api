package cache

import (
	core "fl/my-portfolio/internal/app"
	"fl/my-portfolio/internal/configs"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	config configs.CacherConfig
	*redis.Client
}

func NewRedis() *Redis {
	config := core.GetConfig().Cacher

	opts := &redis.Options{
		Addr: fmt.Sprintf("%s:%s", config.CacherHost, config.CacherPort),
		DB:   config.CacherDB,
	}

	return &Redis{
       config,
       redis.NewClient(opts),
    }
}
