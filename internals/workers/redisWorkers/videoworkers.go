package redisworkers

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/itsmonday/youtube/internals/domain"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddMovieToRedis(ctx context.Context, userid string, videochan chan *domain.Video, errchan chan<- error, db *redis.Client, wg *sync.WaitGroup) {
	defer func() {
		fmt.Println("Video Added in Redis")
		wg.Done()
	}()

	video := <-videochan
	m_video, err := json.Marshal(video)
	if err != nil {
		errchan <- err
		return
	}

	result, err := db.LPush(ctx, "videos:"+userid, string(m_video)).Result()
	if err != nil {
		errchan <- err
		return
	}
	fmt.Printf("\nresult %v\n", result)
	videochan <- video
}

func UpdateMovieInRedis(ctx context.Context, userid string, videoid string, videochan chan *domain.Video, errchan chan<- error, db *redis.Client, wg *sync.WaitGroup) {
	defer func() {
		fmt.Println("Redis Worker completed successfully for Updating Video Details")
		wg.Done()
	}()

	video_from_chan := <-videochan
	fmt.Println("video_from_chan")
	fmt.Println(video_from_chan)
	total_videos_of_user, err := db.LLen(ctx, "videos:"+userid).Result()
	if err != nil {
		errchan <- err
		return
	}

	video_obj_id, _ := primitive.ObjectIDFromHex(videoid)

	fmt.Printf("Total number of videos:%v", total_videos_of_user)
	videos, err := db.LRange(ctx, "videos:"+userid, 0, total_videos_of_user).Result()
	if err != nil {
		errchan <- err
		return
	}

	for idx, video := range videos {
		fmt.Printf("\n Index : %v \n", idx)
		var v domain.Video
		err = json.Unmarshal([]byte(video), &v)
		if err != nil {
			errchan <- err
			return
		}

		videoId := v.ID
		if videoId == video_obj_id {
			fmt.Printf("\n Video ID: %v \n", videoId)
			v.Title = video_from_chan.Title
			v.Description = video_from_chan.Description
			v.Tags = video_from_chan.Tags
			v.ImgURI = video_from_chan.ImgURI

			m_video, err := json.Marshal(v)
			if err != nil {
				errchan <- err
				return
			}
			_, err = db.LSet(ctx, "videos:"+userid, int64(idx), string(m_video)).Result()
			if err != nil {
				errchan <- err
				return
			}
			break
		}
	}

	videochan <- video_from_chan
}

func DeleteMovieInRedis(ctx context.Context, errchan chan<- error, redisdb *redis.Client, wg *sync.WaitGroup, videoid string, userid string) {
	defer func() {
		wg.Done()
	}()

	total_videos, err := redisdb.LLen(ctx, "videos:"+userid).Result()
	if err != nil {
		errchan <- err
		return
	}

	videos, err := redisdb.LRange(ctx, "videos:"+userid, 0, int64(total_videos)).Result()
	if err != nil {
		errchan <- err
		return
	}

	for idx, video := range videos {
		fmt.Printf("\n idx = %v \n", idx)
		var v domain.Video
		if err := json.Unmarshal([]byte(video), &v); err != nil {
			errchan <- err
			return
		}
		if v.ID.Hex() == videoid && v.UserId.Hex() == userid {
			del_res, err := redisdb.LRem(ctx, "videos:"+userid, 1, video).Result()
			if err != nil {
				errchan <- err
			}
			fmt.Printf("\n deleted response %v \n", del_res)
		}
	}
}
