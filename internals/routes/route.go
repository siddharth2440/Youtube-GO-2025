package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/itsmonday/youtube/internals/handlers"
	"github.com/itsmonday/youtube/internals/services"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRouter(redis *redis.Client, mongo *mongo.Client) *gin.Engine {
	route := gin.Default()

	// services
	authService := services.NewAuthService(mongo, redis)
	userservice := services.NewUserService(mongo, redis)

	// handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userservice)

	authRoute := route.Group("/api/v1/auth")
	{
		authRoute.POST("/register", authHandler.UserRegisterHandler)
		authRoute.POST("/login", authHandler.UserLoginHandler)
		authRoute.GET("/logout", authHandler.UserLogoutHandler)
	}

	userRoute := route.Group("/api/v1/user")
	{
		userRoute.PUT("/update-user-details/:userid", userHandler.UpdateUserDetailsHandler)
		userRoute.DELETE("/delete-user/:userid", userHandler.DeleteUserHandler)
	}

	route.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	return route
}
