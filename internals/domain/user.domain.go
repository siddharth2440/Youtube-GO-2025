package domain

import (
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID              primitive.ObjectID `json:"id" bson:"_id"`
	Name            string             `json:"name"`
	Email           string             `json:"email"`
	Password        string             `json:"password"`
	Img             string             `json:"img"`
	Subscribers     int                `json:"subscribers"`
	SubscribedUsers []string           `json:"subscribedUsers"`
	CreatedAt       time.Time          `json:"createdAt"`
	UpdatedAt       time.Time          `json:"UpdatedAt"`
}

func NewUser(user *User) *User {
	return &User{
		ID:              primitive.NewObjectID(),
		Name:            strings.ToLower(user.Name),
		Email:           strings.ToLower(user.Email),
		Password:        strings.ToLower(user.Password),
		Img:             strings.ToLower(user.Img),
		Subscribers:     0,
		SubscribedUsers: make([]string, 0),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}
