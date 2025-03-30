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
		errChan <- fmt.Errorf("user doesnot exists")
		return
	}

	del_res, err := redis.HDel(ctx, "users", "user:"+userid).Result()
	if err != nil {
		errChan <- err
		return
	}
	fmt.Printf("\n delete response from Redis : %v \n", del_res)
}

func GetUserFromRedis(ctx context.Context, userchan chan<- *domain.User, errchan chan<- error, userid string, db *redis.Client) {
	defer func() {
		fmt.Println("Get User from Redis Routine Completed...")
	}()

	var getUser domain.User
	user, err := db.HGet(ctx, "users", "user:"+userid).Result()
	if err != nil {
		println(err.Error())
		errchan <- err
		return
	}

	err = json.Unmarshal([]byte(user), &getUser)
	if err != nil {
		errchan <- err
		return
	}

	fmt.Println(getUser)

	userchan <- &getUser
}

func SubscribeUserInRedis(ctx context.Context, userChan chan *domain.User, errChan chan<- error, db *redis.Client, wg *sync.WaitGroup, userid string, me string) {
	defer func() {
		fmt.Println("User updated in Redis")
		wg.Done()
	}()

	fmt.Printf("\nMy ID: %v\n", me)
	fmt.Printf("\nUserId %v\n", userid)

	my_details := <-userChan
	fmt.Println("my_details")
	fmt.Println(my_details)

	m_details, _ := json.Marshal(my_details)
	redRes, err := db.HGet(ctx, "users", "user:"+me).Result()
	if err != nil {
		errChan <- err
		return
	}

	fmt.Println("After Getting Me")
	fmt.Println(redRes)

	fmt.Println("After Marshal MY details")
	fmt.Println(string(m_details))
	set_details, err := db.HSet(ctx, "users", "user:"+me, string(m_details)).Result()
	if err != nil {
		errChan <- err
		return
	}
	fmt.Println("After Updating Me")
	fmt.Println(set_details)

	var user domain.User

	get_user, err := db.HGet(ctx, "users", "user:"+userid).Result()
	if err != nil {
		errChan <- err
		return
	}
	fmt.Println("After Getting User")
	fmt.Println(get_user)

	err = json.Unmarshal([]byte(get_user), &user)
	if err != nil {
		errChan <- err
		return
	}

	user.Subscribers += 1

	marshaled_user, _ := json.Marshal(user)

	user_updated_res, err := db.HSet(ctx, "users", "user:"+userid, string(marshaled_user)).Result()
	if err != nil {
		errChan <- err
		return
	}
	fmt.Println("After Updating User")
	fmt.Println(user_updated_res)

	userChan <- my_details
}

func UnSubscribeUserInRedis(ctx context.Context, userchan chan *domain.User, errchan chan<- error, userid string, my_id string, db *redis.Client, wg *sync.WaitGroup) {

	defer func() {
		wg.Done()
	}()

	user := <-userchan
	fmt.Printf("\n Useras %v\n", user)
	m_user, _ := json.Marshal(user)

	red_res, err := db.HSet(ctx, "users", "user:"+my_id, string(m_user)).Result()
	if err != nil {
		errchan <- err
		return
	}
	fmt.Printf("\n Redis Result after updating me %v\n", red_res)

	var user_details domain.User
	get_user, err := db.HGet(ctx, "users", "user:"+userid).Result()
	if err != nil {
		errchan <- err
		return
	}

	err = json.Unmarshal([]byte(get_user), &user_details)
	if err != nil {
		errchan <- err
		return
	}

	user_details.Subscribers -= 1
	m_user_details, _ := json.Marshal(user_details)
	update_chan, err := db.HSet(ctx, "users", "user:"+userid, string(m_user_details)).Result()
	if err != nil {
		errchan <- err
		return
	}

	fmt.Printf("\n Redis Result after updating channel %v\n", update_chan)
	userchan <- &user_details
}
