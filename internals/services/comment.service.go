package services

import (
	"context"
	"time"

	"github.com/itsmonday/youtube/internals/domain"
	mongoworkers "github.com/itsmonday/youtube/internals/workers/mongoWorkers"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CommentServiceStruct struct {
	db    *mongo.Client
	redis *redis.Client
}

func NewCommentService(db *mongo.Client, redis *redis.Client) *CommentServiceStruct {
	return &CommentServiceStruct{db, redis}
}

type CommentServiceInterface interface {
	AddCommentService(comment *domain.Comment, userid string, videoid string) (*domain.Comment, error)
	GetCommentsService(videoid string) (*[]domain.Comment, error)
	GetcommentDetails(commentid string) (*domain.Comment, error)
	DeleteCommentService(commentid, userid string) (*domain.Comment, error)
}

func (NCs *CommentServiceStruct) AddCommentService(comment *domain.Comment, userid string, videoid string) (*domain.Comment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	commentchan := make(chan *domain.Comment, 1)
	errchan := make(chan error, 1)

	user_obj_id, _ := primitive.ObjectIDFromHex(userid)
	video_obj_id, _ := primitive.ObjectIDFromHex(videoid)
	newcomment := domain.NewComment(comment, user_obj_id, video_obj_id)

	go mongoworkers.AddComment(ctx, commentchan, errchan, newcomment, NCs.db)
	for {
		select {
		case comment := <-commentchan:
			return comment, nil
		case err := <-errchan:
			return nil, err
		case <-ctx.Done():
			return nil, context.DeadlineExceeded
		}
	}
}

func (NCs *CommentServiceStruct) GetCommentsService(videoid string) (*[]domain.Comment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	commentschan := make(chan *[]domain.Comment, 1)
	errchan := make(chan error, 1)

	go mongoworkers.GetComments(ctx, commentschan, errchan, videoid, NCs.db)
	for {
		select {
		case comments := <-commentschan:
			return comments, nil
		case err := <-errchan:
			return nil, err
		case <-ctx.Done():
			return nil, context.DeadlineExceeded
		}
	}
}

// get comment details
func (NCs *CommentServiceStruct) GetcommentDetails(commentid string) (*domain.Comment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	commentchan := make(chan *domain.Comment, 1)
	errchan := make(chan error, 1)

	// worker to get comment details
	go mongoworkers.GetcommentDetails(ctx, commentchan, errchan, NCs.db, commentid)

	for {
		select {
		case comment := <-commentchan:
			return comment, nil
		case err := <-errchan:
			return nil, err
		case <-ctx.Done():
			return nil, context.DeadlineExceeded
		}
	}

}

// delete comment
func (NCs *CommentServiceStruct) DeleteCommentService(commentid, userid string) (*domain.Comment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	commentchan := make(chan *domain.Comment, 1)
	errorchan := make(chan error, 1)

	// worker to delete the comment
	go mongoworkers.DeleteCommentIMongoDb(ctx, commentchan, errorchan, NCs.db, commentid, userid)

	for {
		select {
		case comment := <-commentchan:
			return comment, nil
		case err := <-errorchan:
			return nil, err
		case <-ctx.Done():
			return nil, context.DeadlineExceeded
		}
	}
}
