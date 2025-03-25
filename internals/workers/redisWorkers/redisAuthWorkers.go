package redisworkers

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/itsmonday/youtube/internals/domain"
	"github.com/redis/go-redis/v9"
)

func RegisterUserRedis(ctx context.Context, userChan chan *domain.User, errChan chan<- error, wg *sync.WaitGroup, redis *redis.Client) {
	defer func() {
		fmt.Println("User Regsitered in redis")
		wg.Done()
	}()

	get_user := <-userChan

	m_user, err := json.Marshal(get_user)
	if err != nil {
		errChan <- fmt.Errorf("error marshalling user: %w", err)
		return
	}

	redis_result, err := redis.HSet(ctx, "users", "user:"+get_user.ID.Hex(), string(m_user)).Result()
	if err != nil {
		errChan <- fmt.Errorf("error setting user in redis: %w", err)
		return
	}
	fmt.Printf("\nredis_result %v\n", redis_result)
	userChan <- get_user
}
