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

func NewComment(comment *Comment, userid primitive.ObjectID, videoid primitive.ObjectID) *Comment {
	return &Comment{
		ID:          primitive.NewObjectID(),
		UserId:      userid,
		VideoId:     videoid,
		Description: strings.ToLower((*comment).Description),
	}
}
