package redisworkers

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/itsmonday/youtube/internals/domain"
	"github.com/redis/go-redis/v9"
)

func UpdateUserDetailsInRedis(ctx context.Context, userid string, userChan chan *domain.User, errChan chan<- error, wg *sync.WaitGroup, redis *redis.Client) {
	defer func() {
		fmt.Println("User Updated Successfully in Redis ")
		defer wg.Done()
	}()

	user := <-userChan

	m_user, err := json.Marshal(user)
	if err != nil {
		errChan <- err
		return
	}

	red_res, err := redis.HSet(ctx, "users", "user:"+userid, string(m_user)).Result()
	if err != nil {
		errChan <- err
		return
	}
	fmt.Printf("\n Redis Result %v\n", red_res)

	userChan <- user
}

func DeleteUserInRedis(ctx context.Context, userid string, errChan chan<- error, wg *sync.WaitGroup, redis *redis.Client) {
	defer func() {
		fmt.Println("Deletion from Redis Completed...")
		wg.Done()
	}()

	isExists, err := redis.HExists(ctx, "users", "user:"+userid).Result()
	if err != nil || !isExists {
		errChan <- fmt.Errorf("User doesnot exists")
		return
	}

	del_res, err := redis.HDel(ctx, "users", "user:"+userid).Result()
	if err != nil {
		errChan <- err
		return
	}
	fmt.Printf("\n delete response from Redis : %v \n", del_res)
}
