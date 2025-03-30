package services

import (
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

type VideoServiceStruct struct {
	db    *mongo.Client
	redis *redis.Client
}

func NewVideoService(db *mongo.Client, redis *redis.Client) *VideoServiceStruct {
	return &VideoServiceStruct{
		db:    db,
		redis: redis,
	}
}

type VideOServiceInterface interface {
}
