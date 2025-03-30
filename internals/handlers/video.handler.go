package handlers

import "github.com/itsmonday/youtube/internals/services"

type VideoHandlerStruct struct {
	service *services.VideOServiceInterface
}

func NewVideoHandler(service *services.VideOServiceInterface) *VideoHandlerStruct {
	return &VideoHandlerStruct{
		service: service,
	}
}
