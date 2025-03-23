package domain

import (
	"fmt"
	"strings"
	"time"

	"github.com/itsmonday/youtube/internals/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID              primitive.ObjectID `bson:"_id" json:"_id"`
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
	id := primitive.NewObjectID()
	fmt.Printf("\nuserid%v\n", id)

	hPassword, _ := utils.HashPassword((*user).Password)
	return &User{
		ID:              primitive.NewObjectID(),
		Name:            strings.ToLower(user.Name),
		Email:           strings.ToLower(user.Email),
		Password:        hPassword,
		Img:             strings.ToLower(user.Img),
		Subscribers:     0,
		SubscribedUsers: make([]string, 0),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func FormatLogin(payload *LoginPayload) *LoginPayload {
	return &LoginPayload{
		Email:    strings.ToLower(payload.Email),
		Password: strings.ToLower(payload.Password),
	}
}
