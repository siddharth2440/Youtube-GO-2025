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

func SubscribeUserInMongoDB(ctx context.Context, userChan chan<- *domain.User, errChan chan<- error, userid string, me string, wg *sync.WaitGroup, db *mongo.Client) {
	defer func() {
		fmt.Println("User Subscribed to a channel in mongodb")
		wg.Done()
	}()

	user_ob_id, _ := primitive.ObjectIDFromHex(userid)
	me_ob_id, _ := primitive.ObjectIDFromHex(me)

	to_find_user := bson.M{
		"_id": user_ob_id,
	}

	to_find_me := bson.M{
		"_id": me_ob_id,
	}

	upda_user, err := db.Database("youtube").Collection("users").UpdateOne(ctx, to_find_user, bson.M{
		"$inc": bson.M{
			"subscribers": 1,
		},
	})
	if err != nil {
		errChan <- err
		return
	}
	fmt.Println("After Updating the Channel Subscribers")
	fmt.Println(upda_user)

	update_me, err := db.Database("youtube").Collection("users").UpdateOne(ctx, to_find_me, bson.M{
		"$addToSet": bson.M{
			"subscribedusers": userid,
		},
	})
	fmt.Println("After Updating the User")
	fmt.Println(update_me)

	if err != nil {
		errChan <- err
		return
	}

	var user domain.User

	if err := db.Database("youtube").Collection("users").FindOne(ctx, to_find_me).Decode(&user); err != nil {
		errChan <- err
		return
	}
	userChan <- &user

}

func UnsubscribeUserInMongoDB(ctx context.Context, userChan chan<- *domain.User, errChan chan<- error, userid string, me string, wg *sync.WaitGroup, db *mongo.Client) {
	defer func() {
		fmt.Println("User Subscribed to a channel in mongodb")
		wg.Done()
	}()

	user_obj_id, _ := primitive.ObjectIDFromHex(userid)
	me_obj_id, _ := primitive.ObjectIDFromHex(me)

	to_find_user := bson.M{
		"_id": user_obj_id,
	}

	to_find_me := bson.M{
		"_id": me_obj_id,
	}

	to_update_user_channel_details := bson.M{
		"$inc": bson.M{
			"subscribers": -1,
		},
	}

	to_update_me := bson.M{
		"$pull": bson.M{
			"subscribedusers": userid,
		},
	}

	chan_update_result, err := db.Database("youtube").Collection("users").UpdateOne(ctx, to_find_user, to_update_user_channel_details)
	if err != nil {
		errChan <- err
		return
	}

	fmt.Printf("\n Chanel Details updated after Unsubscription %v \n", chan_update_result)

	my_details_upated, err := db.Database("youtube").Collection("users").UpdateOne(ctx, to_find_me, to_update_me)
	if err != nil {
		errChan <- err
		return
	}

	fmt.Printf("\n After Unsubscribe to channel my details updated %v \n", my_details_upated)

	var user domain.User

	if err := db.Database("youtube").Collection("users").FindOne(ctx, bson.M{"_id": me_obj_id}).Decode(&user); err != nil {
		errChan <- err
		return
	}
	userChan <- &user
}
