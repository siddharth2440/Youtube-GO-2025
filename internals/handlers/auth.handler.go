package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/itsmonday/youtube/internals/domain"
	"github.com/itsmonday/youtube/internals/services"
)

type AuthHandlerStruct struct {
	service services.AuthServiceInterface
}

func NewAuthHandler(service services.AuthServiceInterface) *AuthHandlerStruct {
	return &AuthHandlerStruct{
		service,
	}
}

// Register
func (NAh *AuthHandlerStruct) UserRegisterHandler(ctx *gin.Context) {
	var user domain.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()})
		return
	}
	userChan := make(chan *domain.User, 1)
	errChan := make(chan error, 1)

	go func(userChan chan<- *domain.User, errChan chan<- error, user *domain.User) {
		user, tokenstring, err := NAh.service.UserRegisterService(user)
		if err != nil {
			errChan <- err
			return
		}
		ctx.SetCookie("youtubecookie", tokenstring, 3600, "/", "localhost", false, true) // 3600 seconds
		userChan <- user
	}(userChan, errChan, &user)

	for {
		select {
		case newUser := <-userChan:
			ctx.JSON(
				http.StatusOK,
				gin.H{
					"success": true,
					"data":    newUser,
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
		}
	}
}

// Login
func (NAh *AuthHandlerStruct) UserLoginHandler(ctx *gin.Context) {
	var loginpayload domain.LoginPayload
	if err := ctx.ShouldBindJSON(&loginpayload); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()})
		return
	}

	userChan := make(chan *domain.User, 1)
	errChan := make(chan error, 1)

	go func(userchan chan<- *domain.User, errChan chan<- error, loginpayload *domain.LoginPayload) {
		loginUser, token, err := NAh.service.UserLoginService(loginpayload)
		if err != nil {
			// fmt.Println("Error  in login handler")
			// fmt.Println(err)
			errChan <- err
			return
		}
		ctx.SetCookie("youtubecookie", token, 3600, "/", "localhost", false, true) // 3600 seconds
		userChan <- loginUser
	}(userChan, errChan, &loginpayload)

	for {
		select {
		case user := <-userChan:
			ctx.JSON(
				http.StatusAccepted,
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

// Logout
func (NAh *AuthHandlerStruct) UserLogoutHandler(ctx *gin.Context) {
	ctx.SetCookie("youtubecookie", "", 0, "/", "localhost", false, true)
	ctx.JSON(
		http.StatusOK,
		gin.H{
			"success": false,
			"message": "Logged out",
		})
}
