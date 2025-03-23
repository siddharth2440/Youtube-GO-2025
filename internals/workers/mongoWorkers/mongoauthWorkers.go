package mongoworkers

import (
	"context"
	"fmt"
	"sync"

	"github.com/itsmonday/youtube/internals/domain"
	"github.com/itsmonday/youtube/internals/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterUserMongoDB(ctx context.Context, user *domain.User, userChan chan<- *domain.User, errChan chan<- error, wg *sync.WaitGroup, db *mongo.Client) {
	defer func() {
		fmt.Println("User Regsitered in MongoDB")
		wg.Done()
	}()

	insert_res, err := db.Database("youtube").Collection("users").InsertOne(ctx, *user)
	if err != nil {
		fmt.Printf("\nerror inserting user %v\n", err)
		errChan <- err
		return
	}

	fmt.Println("Inserted UseriD")
	fmt.Println(insert_res.InsertedID)

	to_get_user := bson.M{
		"_id": user.ID,
	}

	var getUser domain.User
	if err := db.Database("youtube").Collection("users").FindOne(ctx, to_get_user).Decode(&getUser); err != nil {
		errChan <- err
		return
	}
	userChan <- &getUser
}

func LoginUserMongoDB(ctx context.Context, payload *domain.LoginPayload, userChan chan<- *domain.User, errChan chan<- error, db *mongo.Client) {
	fmt.Println(*payload)
	get_user := bson.M{
		"email": (*payload).Email,
	}
	var user domain.User
	if err := db.Database("youtube").Collection("users").FindOne(ctx, get_user).Decode(&user); err != nil {
		fmt.Printf("error in getting or decoding the user: %v", err)
		errChan <- err
		return
	}

	// password verification
	isPasswordCorect, err := utils.VerifyPassword((*payload).Password, user.Password)
	if err != nil || !isPasswordCorect {
		fmt.Println(err)
		errChan <- fmt.Errorf("invalid user credentials")
		return
	}
	fmt.Printf("Error kucch nhi hai")

	userChan <- &user
}
