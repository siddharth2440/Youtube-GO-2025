package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/itsmonday/youtube/internals/handlers"
	"github.com/itsmonday/youtube/internals/middlewares"
	"github.com/itsmonday/youtube/internals/services"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRouter(redis *redis.Client, mongo *mongo.Client) *gin.Engine {
	route := gin.Default()

	// services
	authService := services.NewAuthService(mongo, redis)
	userservice := services.NewUserService(mongo, redis)
	videoService := services.NewVideoService(mongo, redis)

	// handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userservice)
	videoHandler := handlers.NewVideoHandler(videoService)

	authRoute := route.Group("/api/v1/auth")
	{
		authRoute.POST("/register", authHandler.UserRegisterHandler)
		authRoute.POST("/login", authHandler.UserLoginHandler)
		authRoute.GET("/logout", authHandler.UserLogoutHandler)
	}

	userRoute := route.Group("/api/v1/user")
	userRoute.Use(middlewares.AuthMiddleware())
	{
		userRoute.PUT("/update-user-details/:userid", userHandler.UpdateUserDetailsHandler)
		userRoute.DELETE("/delete-user/:userid", userHandler.DeleteUserHandler)
		userRoute.GET("/get-user/:userid", userHandler.GetUserInfoHandler)
		userRoute.GET("/get-users/", userHandler.GetUsersHandler)
		userRoute.GET("/get-user-by-query/", userHandler.GetUsersByQuery)
		userRoute.PUT("/subscribe-user/:userid", userHandler.SubscribeUserHandler)
		userRoute.PUT("/unsubscribe-user/:userid", userHandler.UnsubscribeUserHandler)
	}

	videoRoute := route.Group("/api/v1/video")
	videoRoute.Use(middlewares.AuthMiddleware())
	{
		videoRoute.POST("/add-video", videoHandler.AddVideoHandler)
		videoRoute.PUT("/update-video/:video_id", videoHandler.UpdateVideoHandler)
		videoRoute.DELETE("/delete-video/:video_id", videoHandler.DeleteVideoHandler)
	}

	publicVideoRoute := route.Group("/api/v1/public-video")
	{
		publicVideoRoute.GET("/get-video/:videoid", videoHandler.GetVideoDetailsHandler)
	}

	route.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	return route
}
