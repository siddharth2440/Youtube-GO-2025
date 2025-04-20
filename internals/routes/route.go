package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/itsmonday/youtube/internals/handlers"
	"github.com/itsmonday/youtube/internals/middlewares"
	"github.com/itsmonday/youtube/internals/services"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRouter(redis *redis.Client, mongo *mongo.Client) *gin.Engine {
	route := gin.Default()
	middlewares.Init()

	route.GET("/metrics", gin.WrapH(promhttp.Handler()))
	route.Use(middlewares.TrackMetrics())

	// services
	authService := services.NewAuthService(mongo, redis)
	userservice := services.NewUserService(mongo, redis)
	videoService := services.NewVideoService(mongo, redis)
	commentService := services.NewCommentService(mongo, redis)

	// handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userservice)
	videoHandler := handlers.NewVideoHandler(videoService)
	commentHandler := handlers.NewCommentHandler(commentService)

	authRoute := route.Group("/api/v1/auth")
	authRoute.Use(middlewares.RateLimit())
	{
		authRoute.POST("/register", authHandler.UserRegisterHandler)
		authRoute.POST("/login", authHandler.UserLoginHandler)
		authRoute.GET("/logout", authHandler.UserLogoutHandler)
	}

	userRoute := route.Group("/api/v1/user")
	userRoute.Use(middlewares.AuthMiddleware())
	userRoute.Use(middlewares.RateLimit())
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
	videoRoute.Use(middlewares.RateLimit())
	{
		videoRoute.POST("/add-video", videoHandler.AddVideoHandler)
		videoRoute.PUT("/update-video/:video_id", videoHandler.UpdateVideoHandler)
		videoRoute.DELETE("/delete-video/:video_id", videoHandler.DeleteVideoHandler)
		videoRoute.PUT("/like-video/:videoid", videoHandler.LikeVideoHandler)
		videoRoute.PUT("/dislike-video/:videoid", videoHandler.DislikeVideoHandler)
	}

	publicVideoRoute := route.Group("/api/v1/public-video")
	publicVideoRoute.Use(middlewares.RateLimit())
	{
		publicVideoRoute.GET("/get-video/:videoid", videoHandler.GetVideoDetailsHandler)
		publicVideoRoute.GET("/get-random-videos", videoHandler.GetRandomVideosHandler)
		publicVideoRoute.GET("/search-videos", videoHandler.SearchVideoHandler)
		publicVideoRoute.GET("/trending-videos", videoHandler.TrendingVideosHandler)
	}

	commentProtectedRoute := route.Group("/api/v1/comment")
	commentProtectedRoute.Use(middlewares.AuthMiddleware())
	commentProtectedRoute.Use(middlewares.RateLimit())
	{
		commentProtectedRoute.POST("/add-comment/:videoid", commentHandler.AddCommentHandler)
		commentProtectedRoute.DELETE("/remove-comment/:commentid", commentHandler.DeleteCommentHandler)
	}

	commentpublicroute := route.Group("/api/v1/public-comment")
	commentpublicroute.Use(middlewares.RateLimit())
	{
		commentpublicroute.GET("/get-comments/:videoid", commentHandler.GetCommentsHandler)
		commentpublicroute.GET("/get-comment-details/:commentid", commentHandler.GetCommentDetailsHandler)
	}

	route.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	return route
}
