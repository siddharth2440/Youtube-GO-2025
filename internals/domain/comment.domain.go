package domain

import (
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	UserId      primitive.ObjectID `json:"userid" bson:"userid"`
	VideoId     primitive.ObjectID `json:"videoid" bson:"videoid"`
	Description string             `json:"desc" bson:"desc"`
}

func NewComment(comment *Comment) *Comment {
	return &Comment{
		ID:          primitive.NewObjectID(),
		UserId:      (*comment).UserId,
		VideoId:     (*comment).VideoId,
		Description: strings.ToLower((*comment).Description),
	}
}
