package services

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/itsmonday/youtube/internals/domain"
	"github.com/itsmonday/youtube/internals/utils"
	mongoworkers "github.com/itsmonday/youtube/internals/workers/mongoWorkers"
	redisworkers "github.com/itsmonday/youtube/internals/workers/redisWorkers"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthServiceStruct struct {
	mongo *mongo.Client
	redis *redis.Client
}

type AuthServiceInterface interface {
	UserRegisterService(user *domain.User) (*domain.User, string, error)
	UserLoginService(loginPayload *domain.LoginPayload) (*domain.User, string, error)
}

func NewAuthService(mongo *mongo.Client, redis *redis.Client) *AuthServiceStruct {
	return &AuthServiceStruct{
		mongo: mongo,
		redis: redis,
	}
}

// Register
func (NAs *AuthServiceStruct) UserRegisterService(user *domain.User) (*domain.User, string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	newUser := domain.NewUser(user)

	userChan := make(chan *domain.User, 1)
	errChan := make(chan error, 1)

	var wg sync.WaitGroup
	wg.Add(2)
	// workers

	go mongoworkers.RegisterUserMongoDB(ctx, newUser, userChan, errChan, &wg, NAs.mongo)
	go redisworkers.RegisterUserRedis(ctx, userChan, errChan, &wg, NAs.redis)

	wg.Wait()

	for {
		select {
		case user := <-userChan:
			// Generate Token
			tokenstring, err := utils.GenerateJWTToken(user.ID.Hex(), user.Name, user.Email)
			if err != nil {
				return user, "", nil
			}
			return user, tokenstring, nil
		case err := <-errChan:
			return nil, "", err
		case <-ctx.Done():
			return nil, "", context.DeadlineExceeded
		}
	}
}

// Login
func (NAs *AuthServiceStruct) UserLoginService(loginPayload *domain.LoginPayload) (*domain.User, string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	f_payload := domain.FormatLogin(loginPayload)

	userchan := make(chan *domain.User, 1)
	errchan := make(chan error, 1)
	go mongoworkers.LoginUserMongoDB(ctx, f_payload, userchan, errchan, NAs.mongo)

	for {
		select {
		case user := <-userchan:
			token, err := utils.GenerateJWTToken(user.ID.Hex(), user.Name, user.Email)
			if err != nil {
				return nil, "", err
			}
			return user, token, nil
		case err := <-errchan:
			fmt.Printf("\n Error coming in channel%v\n", err)
			fmt.Println(err)
			return nil, "", err
		case <-ctx.Done():
			return nil, "", context.DeadlineExceeded
		}
	}
}
