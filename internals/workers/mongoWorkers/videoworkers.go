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

func AddAVideoInMongoDB(ctx context.Context, video *domain.Video, videochan chan<- *domain.Video, errchan chan error, db *mongo.Client, wg *sync.WaitGroup) {
	defer func() {
		fmt.Println("Video Added inside the MongoDB")
		wg.Done()
	}()

	inserted_res, err := db.Database("youtube").Collection("videos").InsertOne(ctx, *video)
	if err != nil {
		errchan <- err
		return
	}

	fmt.Printf("\nINserted Result %v \n", inserted_res)

	video_id := (*video).ID

	to_get_video := bson.M{
		"_id": video_id,
	}

	var get_video domain.Video
	if err := db.Database("youtube").Collection("videos").FindOne(ctx, to_get_video).Decode(&get_video); err != nil {
		errchan <- err
		return
	}

	videochan <- &get_video
}

func UpdateAVideoInMongoDB(ctx context.Context, videoChan chan<- *domain.Video, updatePayload *domain.UpdateVideoPayload, db *mongo.Client, errChan chan<- error, userid string, videoId string, wg *sync.WaitGroup) {
	defer func() {
		fmt.Println("Video Updated")
		wg.Done()
	}()

	video_obj_id, err := primitive.ObjectIDFromHex(videoId)
	if err != nil {
		errChan <- err
		return
	}
	user_obj_id, err := primitive.ObjectIDFromHex(userid)
	if err != nil {
		errChan <- err
		return
	}

	to_find_video := bson.M{
		"$and": bson.A{
			bson.M{
				"_id": video_obj_id,
			},
			bson.M{
				"user_id": user_obj_id,
			},
		},
	}
	var video domain.Video

	if err := db.Database("youtube").Collection("videos").FindOne(ctx, to_find_video).Decode(&video); err != nil {
		errChan <- err
		return
	}

	// update video Details
	video.Title = updatePayload.Title
	video.Description = updatePayload.Description
	video.Tags = updatePayload.Tags
	video.ImgURI = updatePayload.ImgURI

	fmt.Println(video)
	to_update := bson.M{
		"$set": video,
	}

	update_res, err := db.Database("youtube").Collection("videos").UpdateOne(ctx, to_find_video, to_update)
	if err != nil {
		errChan <- err
		return
	}
	fmt.Printf("\n After Updating a video %v \n", update_res)
	videoChan <- &video
}

func DeleteVideoFromMongoDB(ctx context.Context, userid string, videoid string, videochan chan<- *domain.Video, errchan chan<- error, db *mongo.Client, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()

	user_obj_id, _ := primitive.ObjectIDFromHex(userid)
	video_obj_id, _ := primitive.ObjectIDFromHex(videoid)

	to_find_video := bson.M{
		"$and": bson.A{
			bson.M{
				"_id": video_obj_id,
			},
			bson.M{
				"user_id": user_obj_id,
			},
		},
	}
	var video domain.Video
	if err := db.Database("youtube").Collection("videos").FindOne(ctx, to_find_video).Decode(&video); err != nil {
		errchan <- err
		return
	}

	del_res, err := db.Database("youtube").Collection("videos").DeleteOne(ctx, to_find_video)
	if err != nil {
		errchan <- err
		return
	}
	fmt.Printf("\n Video Deleted Result %v \n", del_res)
	videochan <- &video
}

func GetMovieDetails(ctx context.Context, videochan chan<- *domain.Video, errchan chan<- error, videoid string, db *mongo.Client) {
	video_obj_id, _ := primitive.ObjectIDFromHex(videoid)
	to_find_video := bson.M{
		"_id": video_obj_id,
	}

	var video domain.Video

	if err := db.Database("youtube").Collection("videos").FindOne(ctx, to_find_video).Decode(&video); err != nil {
		errchan <- err
		return
	}
	videochan <- &video
}

func RandomVideos(ctx context.Context, videochan chan<- *[]domain.Video, errchan chan<- error, db *mongo.Client, limit string) {

	c_limit, _ := strconv.Atoi(limit)
	to_get_random_videos := bson.A{
		bson.M{
			"$sample": bson.M{
				"size": c_limit,
			},
		},
	}

	var videos []domain.Video
	cur, err := db.Database("youtube").Collection("videos").Aggregate(ctx, to_get_random_videos)
	if err != nil {
		errchan <- err
		return
	}

	for cur.Next(ctx) {
		var video domain.Video
		if err := cur.Decode(&video); err != nil {
			errchan <- err
			return
		}
		videos = append(videos, video)
	}
	if err := cur.Close(ctx); err != nil {
		errchan <- err
		return
	}
	videochan <- &videos
}

func SearchVideoInMongodb(ctx context.Context, query string, videoschan chan<- *[]domain.Video, errchan chan<- error, db *mongo.Client) {

	to_search_videos := bson.M{
		"$or": bson.A{
			bson.M{
				"title": bson.M{
					"$regex":   query,
					"$options": "i",
				},
			},
			bson.M{
				"description": bson.M{
					"$regex":   query,
					"$options": "i",
				},
			},
			bson.M{
				"tags": bson.M{
					"$in": []string{query},
				},
			},
		},
		"user_id": bson.M{
			"$exists": true,
		},
	}

	cur, err := db.Database("youtube").Collection("videos").Find(ctx, to_search_videos)
	if err != nil {
		errchan <- err
		return
	}

	var videos []domain.Video
	for cur.Next(ctx) {
		var video domain.Video
		if err := cur.Decode(&video); err != nil {
			errchan <- err
			return
		}
		videos = append(videos, video)
	}
	if err := cur.Close(ctx); err != nil {
		errchan <- err
		return
	}
	videoschan <- &videos
}

func GetTrendingVideosFromMongodb(ctx context.Context, videoschan chan<- *[]domain.Video, errchan chan<- error, db *mongo.Client) {
	to_get_trending_videos := bson.A{
		bson.M{
			"$sort": bson.M{
				"views": -1,
			},
		},
		bson.M{
			"$limit": 2,
		},
	}

	cur, err := db.Database("youtube").Collection("videos").Aggregate(ctx, to_get_trending_videos)
	if err != nil {
		errchan <- err
		return
	}

	var videos []domain.Video

	for cur.Next(ctx) {
		var video domain.Video
		if err := cur.Decode(&video); err != nil {
			errchan <- err
			return
		}
		videos = append(videos, video)
	}
	videoschan <- &videos
}

func LikeVideo(ctx context.Context, videochan chan<- *domain.Video, errchan chan<- error, videoid string, userid string, db *mongo.Client) {
	user_obj_id, _ := primitive.ObjectIDFromHex(userid)
	video_obj_id, _ := primitive.ObjectIDFromHex(videoid)

	to_find_video := bson.M{
		"_id": video_obj_id,
	}

	update_video := bson.M{
		"$addToSet": bson.M{
			"likes": user_obj_id,
		},
		"$pull": bson.M{
			"dislikes": user_obj_id,
		},
	}

	update_result, err := db.Database("youtube").Collection("videos").UpdateOne(ctx, to_find_video, update_video)

	if err != nil {
		errchan <- err
		return
	}

	fmt.Println("Updated Result")
	fmt.Println(update_result)

	var liked_video domain.Video
	if err := db.Database("youtube").Collection("videos").FindOne(ctx, to_find_video).Decode(&liked_video); err != nil {
		errchan <- err
		return
	}
	videochan <- &liked_video
}

func DislikeVideo(ctx context.Context, videochan chan<- *domain.Video, errchan chan<- error, videoid string, userid string, db *mongo.Client) {
	user_obj_id, _ := primitive.ObjectIDFromHex(userid)
	video_obj_id, _ := primitive.ObjectIDFromHex(videoid)

	to_find_video := bson.M{
		"_id": video_obj_id,
	}

	update_video := bson.M{
		"$addToSet": bson.M{
			"dislikes": user_obj_id,
		},
		"$pull": bson.M{
			"likes": user_obj_id,
		},
	}

	updated_video_res, err := db.Database("youtube").Collection("videos").UpdateOne(ctx, to_find_video, update_video)
	if err != nil {
		errchan <- err
		return
	}
	fmt.Println("Updated Result")
	fmt.Println(updated_video_res)

	var disliked_video domain.Video
	if err := db.Database("youtube").Collection("videos").FindOne(ctx, to_find_video).Decode(&disliked_video); err != nil {
		errchan <- err
		return
	}
	videochan <- &disliked_video
}
