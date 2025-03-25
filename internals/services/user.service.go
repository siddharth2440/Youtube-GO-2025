package services

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/itsmonday/youtube/internals/domain"
	mongoworkers "github.com/itsmonday/youtube/internals/workers/mongoWorkers"
	redisworkers "github.com/itsmonday/youtube/internals/workers/redisWorkers"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserServiceInterface interface {
	UpdateUserDetailsService(userDetails *domain.UpdateDetails, userid string) (*domain.User, error)
	DeleteUserService(userid string) (*domain.User, error)
}

type UserServiceStruct struct {
	db    *mongo.Client
	redis *redis.Client
}

func NewUserService(db *mongo.Client, redis *redis.Client) *UserServiceStruct {
	return &UserServiceStruct{db, redis}
}

func (NUs *UserServiceStruct) UpdateUserDetailsService(userDetails *domain.UpdateDetails, userid string) (*domain.User, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	userChan := make(chan *domain.User, 1)
	errChan := make(chan error, 1)

	userDetails = domain.NewUpdateDetails(userDetails)
	fmt.Println(userDetails)

	var wg sync.WaitGroup
	wg.Add(2)

	// Mongodb
	go mongoworkers.UpdateUserDetailsInMongoDB(userDetails, ctx, userid, NUs.db, userChan, errChan, &wg)
	// Redis
	go redisworkers.UpdateUserDetailsInRedis(ctx, userid, userChan, errChan, &wg, NUs.redis)
	wg.Wait()

	for {
		select {
		case user := <-userChan:
			return user, nil
		case err := <-errChan:
			return nil, err
		case <-ctx.Done():
			return nil, context.DeadlineExceeded
		}
	}
}

func (NUs *UserServiceStruct) DeleteUserService(userid string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	userChan := make(chan *domain.User, 1)
	errChan := make(chan error, 1)

	var wg sync.WaitGroup
	wg.Add(2)

	// mongodb go routine
	go mongoworkers.DeleteUserFromMongoDB(userid, ctx, userChan, errChan, &wg, NUs.db)
	// redis routine
	go redisworkers.DeleteUserInRedis(ctx, userid, errChan, &wg, NUs.redis)

	wg.Wait()

	for {
		select {
		case user := <-userChan:
			return user, nil
		case err := <-errChan:
			return nil, err
		case <-ctx.Done():
			return nil, context.DeadlineExceeded
		}
	}
}
