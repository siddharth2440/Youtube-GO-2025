package mongoworkers

import (
	"context"

	"github.com/itsmonday/youtube/internals/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddComment(ctx context.Context, commentchan chan<- *domain.Comment, errchan chan<- error, comment *domain.Comment, db *mongo.Client) {
	insert_res, err := db.Database("youtube").Collection("comments").InsertOne(ctx, (*comment))
	if err != nil {
		errchan <- err
		return
	}

	insert_id := insert_res.InsertedID

	to_get_comment := bson.M{
		"_id": insert_id,
	}
	var get_comment domain.Comment
	if err := db.Database("youtube").Collection("comments").FindOne(ctx, to_get_comment).Decode(&get_comment); err != nil {
		errchan <- err
		return
	}
	commentchan <- &get_comment
}

func GetComments(ctx context.Context, commentschan chan<- *[]domain.Comment, errchan chan<- error, videoid string, db *mongo.Client) {
	video_obj_id, err := primitive.ObjectIDFromHex(videoid)
	if err != nil {
		errchan <- err
		return
	}
	to_get_comments := bson.M{
		"videoid": video_obj_id,
	}
	cur, err := db.Database("youtube").Collection("comments").Find(ctx, to_get_comments)
	if err != nil {
		errchan <- err
		return
	}
	var comments []domain.Comment
	for cur.Next(ctx) {
		var comment domain.Comment
		if err := cur.Decode(&comment); err != nil {
			errchan <- err
			return
		}
		comments = append(comments, comment)
	}

	defer func() {
		if err := cur.Close(ctx); err != nil {
			errchan <- err
			return
		}
	}()

	commentschan <- &comments
}

func GetcommentDetails(ctx context.Context, commentchan chan<- *domain.Comment, errchan chan<- error, db *mongo.Client, commentid string) {
	comment_obj_id, err := primitive.ObjectIDFromHex(commentid)
	if err != nil {
		errchan <- err
		return
	}

	to_find_comment := bson.M{
		"_id": comment_obj_id,
	}
	var comment domain.Comment
	if err := db.Database("youtube").Collection("comments").FindOne(ctx, to_find_comment).Decode(&comment); err != nil {
		errchan <- err
		return
	}
	commentchan <- &comment
}

func DeleteCommentIMongoDb(ctx context.Context, commentchan chan<- *domain.Comment, errchan chan<- error, db *mongo.Client, commentid string, userid string) {
	comment_obj_id, err := primitive.ObjectIDFromHex(commentid)
	if err != nil {
		errchan <- err
		return
	}

	user_obj_id, err := primitive.ObjectIDFromHex(userid)
	if err != nil {
		errchan <- err
		return
	}

	to_find_comment := bson.M{
		"$and": bson.A{
			bson.M{
				"_id": comment_obj_id,
			},
			bson.M{
				"userid": user_obj_id,
			},
		},
	}

	var deleted_comment_details domain.Comment
	if err := db.Database("youtube").Collection("comments").FindOneAndDelete(ctx, to_find_comment).Decode(&deleted_comment_details); err != nil {
		errchan <- err
		return
	}
	commentchan <- &deleted_comment_details
}
