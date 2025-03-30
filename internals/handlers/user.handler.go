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
func (NUh *UserHandlerStruct) GetUserInfoHandler(ctx *gin.Context) {
	userId := ctx.Param("userid")

	userchan := make(chan *domain.User, 1)
	errchan := make(chan error, 1)

	go func(userchan chan<- *domain.User, errchan chan<- error, userId string) {
		user, err := NUh.services.GetUserService(userId)
		if err != nil {
			errchan <- err
			return
		}
		userchan <- user
	}(userchan, errchan, userId)

	for {
		select {
		case user := <-userchan:
			ctx.JSON(
				http.StatusOK,
				gin.H{
					"success": true,
					"data":    user,
				},
			)
			return
		case err := <-errchan:
			ctx.JSON(
				http.StatusInternalServerError,
				gin.H{
					"success": false,
					"error":   err.Error(),
				},
			)
			return
		}
	}
}

// get users
func (NUh *UserHandlerStruct) GetUsersHandler(ctx *gin.Context) {
	query := ctx.Query("users")
	userchan := make(chan *[]domain.User, 1)
	errchan := make(chan error, 1)

	go func(userchan chan<- *[]domain.User, errchan chan<- error, query string) {
		users, err := NUh.services.GetUsersService(query)
		if err != nil {
			errchan <- err
			return
		}
		userchan <- users
	}(userchan, errchan, query)
	for {
		select {
		case user := <-userchan:
			ctx.JSON(
				http.StatusOK,
				gin.H{
					"success": true,
					"data":    user,
				},
			)
			return
		case err := <-errchan:
			ctx.JSON(
				http.StatusInternalServerError,
				gin.H{
					"success": false,
					"error":   err.Error(),
				},
			)
			return
		}
	}
}

// find by query := "email" or "username"
func (NUh *UserHandlerStruct) GetUsersByQuery(ctx *gin.Context) {
	query := ctx.Query("query")

	userchan := make(chan *[]domain.User, 1)
	errchan := make(chan error, 1)

	go func(userchan chan<- *[]domain.User, errchan chan<- error, query string) {
		usrs, err := NUh.services.GetUserByQuery(query)
		if err != nil {
			errchan <- err
			return
		}
		userchan <- usrs
	}(userchan, errchan, query)

	for {
		select {
		case user := <-userchan:
			ctx.JSON(
				http.StatusOK,
				gin.H{
					"success": true,
					"data":    user,
				},
			)
			return
		case err := <-errchan:
			ctx.JSON(
				http.StatusInternalServerError,
				gin.H{
					"success": false,
					"error":   err.Error(),
				},
			)
			return
		}
	}
}

func (NUh *UserHandlerStruct) SubscribeUserHandler(ctx *gin.Context) {
	userId := ctx.Param("userid")
	me := ctx.GetString("authuserid")

	fmt.Printf("\n  userId := %v\n", userId)
	fmt.Printf("\n  me := %v\n", me)

	userChan := make(chan *domain.User, 1)
	errChan := make(chan error, 1)

	go func(userchan chan<- *domain.User, errchan chan<- error, userid string, me string) {
		user, err := NUh.services.SubscribeUserService(userid, me)
		if err != nil {
			errChan <- err
			return
		}
		userChan <- user

	}(userChan, errChan, userId, me)
	for {
		select {
		case user := <-userChan:
			ctx.JSON(
				http.StatusOK,
				gin.H{
					"success": true,
					"data":    user,
				},
			)
			return
		case err := <-errChan:
			ctx.JSON(
				http.StatusInternalServerError,
				gin.H{
					"success": false,
					"error":   err.Error(),
				},
			)
			return
		}
	}
}

func (NUh *UserHandlerStruct) UnsubscribeUserHandler(ctx *gin.Context) {
	userid := ctx.Param("userid")
	me := ctx.GetString("authuserid")

	userChan := make(chan *domain.User, 1)
	errChan := make(chan error, 1)

	go func(userChan chan<- *domain.User, errChan chan<- error, userid string, me string) {
		user, err := NUh.services.UnsubscribeUserService(userid, me)
		if err != nil {
			errChan <- err
			return
		}
		userChan <- user
	}(userChan, errChan, userid, me)

	for {
		select {
		case user := <-userChan:
			ctx.JSON(
				http.StatusOK,
				gin.H{
					"success": true,
					"data":    user,
				},
			)
			return

		case err := <-errChan:
			ctx.JSON(
				http.StatusInternalServerError,
				gin.H{
					"success": true,
					"error":   err.Error(),
				},
			)
			return
		}
	}

}
