package mongoworkers

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/itsmonday/youtube/internals/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func UpdateUserDetailsInMongoDB(updateDetails *domain.UpdateDetails, ctx context.Context, userid string, db *mongo.Client, userChan chan<- *domain.User, errChan chan<- error, wg *sync.WaitGroup) {
	defer func() {
		fmt.Println("User Updated Successfully in MongoDB ")
		defer wg.Done()
	}()

	// string to objectID
	u_id, err := primitive.ObjectIDFromHex(userid)
	if err != nil {
		errChan <- fmt.Errorf("error in the conversion from string to objectId: %v", err)
		return
	}

	to_get_user := bson.M{
		"_id": u_id,
	}

	var get_user domain.User
	if err := db.Database("youtube").Collection("users").FindOne(ctx, to_get_user).Decode(&get_user); err != nil {
		errChan <- err
		return
	}

	// update information
	get_user.Email = updateDetails.Email
	get_user.Name = updateDetails.Name

	to_update_details := bson.M{
		"$set": get_user,
	}

	update_result, err := db.Database("youtube").Collection("users").UpdateOne(ctx, to_get_user, to_update_details)
	if err != nil {
		errChan <- err
		return
	}

	fmt.Printf("\nupdate_result %v\n", update_result)
	userChan <- &get_user
}

func DeleteUserFromMongoDB(userid string, ctx context.Context, userchan chan<- *domain.User, errchan chan<- error, wg *sync.WaitGroup, db *mongo.Client) {
	defer func() {
		fmt.Println("User Deletion in Mongodb completed successfully")
		wg.Done()
	}()

	user_obj_id, err := primitive.ObjectIDFromHex(userid)
	if err != nil {
		errchan <- fmt.Errorf("error converting string to object id: %v", err)
		return
	}

	to_get_user := bson.M{
		"_id": user_obj_id,
	}

	var user domain.User
	if err := db.Database("youtube").Collection("users").FindOne(ctx, to_get_user).Decode(&user); err != nil {
		errchan <- err
		return
	}

	del_res, err := db.Database("youtube").Collection("users").DeleteOne(ctx, to_get_user)
	if err != nil {
		errchan <- err
		return
	}
	fmt.Printf("\ndelete_res %v\n", del_res)
	userchan <- &user
}

func GetUsers(ctx context.Context, userChan chan<- *[]domain.User, errChan chan<- error, db *mongo.Client, query string) {
	defer func() {
		fmt.Println("Get Users Routine Completed...")
	}()

	// query := retrive number of users, e.g., user=10
	users, _ := strconv.Atoi(query)
	fmt.Printf("\n no. of users %v\n", users)
	fmt.Printf("\n no. of users (query) %v\n", query)

	to_get_users := bson.A{
		bson.M{
			"$limit": users,
		},
		bson.M{
			"$sort": bson.M{
				"createdAt": -1,
			},
		},
	}

	cur, err := db.Database("youtube").Collection("users").Aggregate(ctx, to_get_users)
	if err != nil {
		errChan <- err
		return
	}

	var get_users []domain.User
	for cur.Next(ctx) {
		var user domain.User
		err := cur.Decode(&user)
		if err != nil {
			errChan <- err
			return
		}
		get_users = append(get_users, user)
	}

	defer cur.Close(ctx)

	userChan <- &get_users
}

func GetUserByQuery(ctx context.Context, userChan chan<- *[]domain.User, query string, errChan chan<- error, db *mongo.Client) {

	to_get_user_by_query := bson.M{
		"$or": bson.A{
			bson.M{
				"email": bson.M{
					"$regex":   query,
					"$options": "i",
				},
			},

			bson.M{
				"name": bson.M{
					"$regex":   query,
					"$options": "i",
				},
			},
		},
	}

	cur, err := db.Database("youtube").Collection("users").Find(ctx, to_get_user_by_query)
	if err != nil {
		errChan <- err
		return
	}

	var get_users []domain.User
	for cur.Next(ctx) {
		var user domain.User
		err := cur.Decode(&user)
		if err != nil {
			errChan <- err
			return
		}
		get_users = append(get_users, user)
	}

	userChan <- &get_users
}
