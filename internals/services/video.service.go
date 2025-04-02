package services

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/itsmonday/youtube/internals/domain"
	mongoworkers "github.com/itsmonday/youtube/internals/workers/mongoWorkers"
	redisworkers "github.com/itsmonday/youtube/internals/workers/redisWorkers"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type VideoServiceStruct struct {
	db    *mongo.Client
	redis *redis.Client
}

func NewVideoService(db *mongo.Client, redis *redis.Client) *VideoServiceStruct {
	return &VideoServiceStruct{
		db:    db,
		redis: redis,
	}
}

type VideOServiceInterface interface {
	AddAVideoService(video *domain.Video, userid string) (*domain.Video, error)
	UpdateVideoService(payload *domain.UpdateVideoPayload, userid string, videoid string) (*domain.Video, error)
	DeleteVideoService(userid, videoid string) (*domain.Video, error)
	GetVideoDetailsService(videoid string) (*domain.Video, error)
}

// Add a Video
func (NVs *VideoServiceStruct) AddAVideoService(video *domain.Video, userid string) (*domain.Video, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	videoChan := make(chan *domain.Video, 1)
	errChan := make(chan error, 1)

	user_obj_id, err := primitive.ObjectIDFromHex(userid)
	if err != nil {
		errChan <- fmt.Errorf("error converting string to object id: %v", err)
	}
	newVideo := domain.NewVideo(video, user_obj_id)

	var wg sync.WaitGroup
	wg.Add(2)

	// Mongo Worker
	go mongoworkers.AddAVideoInMongoDB(ctx, newVideo, videoChan, errChan, NVs.db, &wg)
	// Redis Worker
	go redisworkers.AddMovieToRedis(ctx, userid, videoChan, errChan, NVs.redis, &wg)

	wg.Wait()

	for {
		select {
		case video := <-videoChan:
			return video, nil
		case err := <-errChan:
			return nil, err
		case <-ctx.Done():
			return nil, context.DeadlineExceeded
		}
	}
}

// Update a Video
func (NVs *VideoServiceStruct) UpdateVideoService(payload *domain.UpdateVideoPayload, userid string, videoid string) (*domain.Video, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	videoChan := make(chan *domain.Video, 1)
	errChan := make(chan error, 1)

	updatePayload := domain.NewUpdatePayload(payload)
	var wg sync.WaitGroup
	wg.Add(2)

	go mongoworkers.UpdateAVideoInMongoDB(ctx, videoChan, updatePayload, NVs.db, errChan, userid, videoid, &wg)
	go redisworkers.UpdateMovieInRedis(ctx, userid, videoid, videoChan, errChan, NVs.redis, &wg)

	wg.Wait()
	for {
		select {
		case video := <-videoChan:
			return video, nil
		case err := <-errChan:
			return nil, err
		case <-ctx.Done():
			return nil, context.DeadlineExceeded
		}
	}
}

// Delete A video
func (NVs *VideoServiceStruct) DeleteVideoService(userid, videoid string) (*domain.Video, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	videochan := make(chan *domain.Video, 1)
	errchan := make(chan error, 1)

	var wg sync.WaitGroup
	wg.Add(2)

	go mongoworkers.DeleteVideoFromMongoDB(ctx, userid, videoid, videochan, errchan, NVs.db, &wg)
	go redisworkers.DeleteMovieInRedis(ctx, errchan, NVs.redis, &wg, videoid, userid)

	wg.Wait()
	for {
		select {
		case video := <-videochan:
			return video, nil
		case err := <-errchan:
			return nil, err
		case <-ctx.Done():
			return nil, context.DeadlineExceeded
		}
	}
}

func (NVs *VideoServiceStruct) GetVideoDetailsService(videoid string) (*domain.Video, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	videochan := make(chan *domain.Video, 1)
	errorchan := make(chan error, 1)

	go mongoworkers.GetMovieDetails(ctx, videochan, errorchan, videoid, NVs.db)

	for {
		select {
		case video := <-videochan:
			return video, nil
		case err := <-errorchan:
			return nil, err
		case <-ctx.Done():
			return nil, context.DeadlineExceeded
		}
	}
}
