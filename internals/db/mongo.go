package db

import (
	"context"

	"github.com/itsmonday/youtube/configs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetMongo(cfg *configs.Env) (*mongo.Client, error) {
	mongoUrl := cfg.MONGO_URL
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	optns := options.Client().ApplyURI(mongoUrl).SetServerAPIOptions(serverAPI)

	mongoClient, err := mongo.Connect(context.TODO(), optns)
	if err != nil {
		return nil, err
	}
	err = mongoClient.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	return mongoClient, nil
}
