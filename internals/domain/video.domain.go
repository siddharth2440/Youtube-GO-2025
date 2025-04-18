package domain

import (
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Video struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	ImgURI      string             `json:"img_uri" bson:"img_uri"`
	VideoURI    string             `json:"video_uri" bson:"video_uri"`
	Views       int                `json:"views" bson:"views"`
	Tags        []string           `json:"tags" bson:"tags"`
	Likes       []string           `json:"likes" bson:"likes"`
	Dislikes    []string           `json:"dislikes" bson:"dislikes"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt"`
	UserId      primitive.ObjectID `json:"user_id" bson:"user_id"`
}

func NewVideo(video *Video, userid primitive.ObjectID) *Video {
	return &Video{
		ID:          primitive.NewObjectID(),
		UserId:      userid,
		Title:       strings.ToLower((*video).Title),
		Description: strings.ToLower((*video).Description),
		ImgURI:      strings.ToLower((*video).ImgURI),
		VideoURI:    strings.ToLower((*video).VideoURI),
		Views:       0,
		Tags:        make([]string, 0),
		Likes:       make([]string, 0),
		Dislikes:    make([]string, 0),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

type UpdateVideoPayload struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	ImgURI      string   `json:"img_uri"`
}

func NewUpdatePayload(payload *UpdateVideoPayload) *UpdateVideoPayload {
	return &UpdateVideoPayload{
		Title:       (*payload).Title,
		Description: (*payload).Description,
		Tags:        (*payload).Tags,
		ImgURI:      (*payload).ImgURI,
	}
}
