package main

import (
	"log"

	"github.com/itsmonday/youtube/configs"
	"github.com/itsmonday/youtube/internals/db"
	"github.com/itsmonday/youtube/internals/routes"
)

func main() {
	// setup configuration
	env, err := configs.Config()

	if err != nil {
		log.Fatalf("Error in setup the config: %v", err)
		return
	}

	// mongodb initialization
	mongoClient, err := db.GetMongo(env)
	if err != nil {
		log.Fatalf("Error in setup the MongoDB: %v", err)
		return
	}

	// redis initialization
	redisClient, err := db.Redis(env)
	if err != nil {
		log.Fatalf("Error in setup the Redis: %v", err)
		return
	}

	server := routes.SetupRouter(redisClient, mongoClient)
	server.Run(":7565")
}
