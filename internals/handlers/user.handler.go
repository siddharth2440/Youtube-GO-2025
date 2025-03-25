package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/itsmonday/youtube/internals/domain"
	"github.com/itsmonday/youtube/internals/services"
)

type UserHandlerStruct struct {
	services services.UserServiceInterface
}

func NewUserHandler(service services.UserServiceInterface) *UserHandlerStruct {
	return &UserHandlerStruct{services: service}
}

// update details
func (NUh *UserHandlerStruct) UpdateUserDetailsHandler(ctx *gin.Context) {
	var u_user domain.UpdateDetails
	userId := ctx.Param("userid")
	fmt.Println(userId)

	if err := ctx.ShouldBindJSON(&u_user); err != nil {
		fmt.Println(err)
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"success": true,
				"error":   fmt.Errorf("%v", err),
			},
		)
		return
	}

	userChan := make(chan *domain.User, 1)
	errChan := make(chan error, 1)

	go func(userchan chan<- *domain.User, errChan chan<- error, user *domain.UpdateDetails, userId string) {
		updatedUser, err := NUh.services.UpdateUserDetailsService(user, userId)
		if err != nil {
			errChan <- err
			return
		}
		userchan <- updatedUser
	}(userChan, errChan, &u_user, userId)

	for {
		select {
		case user := <-userChan:
			ctx.JSON(http.StatusOK, gin.H{
				"success": true,
				"data":    user,
			})
			return
		case err := <-errChan:
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}
	}

}

// delete user
func (NUh *UserHandlerStruct) DeleteUserHandler(ctx *gin.Context) {
	userid := ctx.Param("userid")

	userChan := make(chan *domain.User, 1)
	errChan := make(chan error, 1)

	go func(userChan chan<- *domain.User, errChan chan<- error, userid string) {
		user, err := NUh.services.DeleteUserService(userid)
		if err != nil {
			errChan <- err
			return
		}
		userChan <- user
	}(userChan, errChan, userid)
	for {
		select {
		case user := <-userChan:
			ctx.JSON(http.StatusOK, gin.H{
				"success": true,
				"data":    user,
			})
			return
		case err := <-errChan:
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}
	}
}

// get user info
// get users
// find by query
