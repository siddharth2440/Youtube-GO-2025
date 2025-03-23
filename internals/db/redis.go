package db

import (
	"github.com/itsmonday/youtube/configs"
	"github.com/redis/go-redis/v9"
)

func Redis(cfg *configs.Env) (*redis.Client, error) {
	redisOptions, err := redis.ParseURL(cfg.REDIS_URL)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(redisOptions)
	return client, nil
}
