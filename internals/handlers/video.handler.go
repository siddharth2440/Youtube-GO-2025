package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/itsmonday/youtube/internals/domain"
	"github.com/itsmonday/youtube/internals/services"
)

type VideoHandlerStruct struct {
	service services.VideOServiceInterface
}

func NewVideoHandler(service services.VideOServiceInterface) *VideoHandlerStruct {
	return &VideoHandlerStruct{
		service: service,
	}
}

// Add Video
func (NVh *VideoHandlerStruct) AddVideoHandler(ctx *gin.Context) {
	var video domain.Video
	if err := ctx.ShouldBindJSON(&video); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"success": false,
				"error":   err.Error(),
			},
		)
	}

	userid := ctx.GetString("authuserid")

	videochan := make(chan *domain.Video, 1)
	errchan := make(chan error, 1)

	go func(videochan chan<- *domain.Video, errchan chan<- error, video *domain.Video, userid *string) {
		newvideo, err := NVh.service.AddAVideoService(video, *userid)
		if err != nil {
			errchan <- err
			return
		}
		videochan <- newvideo
	}(videochan, errchan, &video, &userid)
	for {
		select {
		case video := <-videochan:
			ctx.JSON(http.StatusCreated, gin.H{
				"success": true,
				"video":   video,
			})
			return
		case err := <-errchan:
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}
	}
}

// Update Video
func (NVh *VideoHandlerStruct) UpdateVideoHandler(ctx *gin.Context) {
	var video_update_payload domain.UpdateVideoPayload
	if err := ctx.ShouldBindJSON(&video_update_payload); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"success": false,
				"error":   err.Error(),
			},
		)
	}

	userid := ctx.GetString("authuserid")
	videoid := ctx.Param("video_id")
	update_video_chan := make(chan *domain.Video, 1)
	errchan := make(chan error, 1)

	go func(update_video_chan chan<- *domain.Video, errchan chan<- error, videoid *string, userid *string, payload *domain.UpdateVideoPayload) {
		updated_video, err := NVh.service.UpdateVideoService(payload, *userid, *videoid)
		if err != nil {
			errchan <- err
			return
		}

		update_video_chan <- updated_video
	}(update_video_chan, errchan, &videoid, &userid, &video_update_payload)

	for {
		select {
		case video := <-update_video_chan:
			ctx.JSON(http.StatusOK, gin.H{
				"success": true,
				"video":   video,
			})
			return
		case err := <-errchan:
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}
	}
}

// Delete Video
func (NVh *VideoHandlerStruct) DeleteVideoHandler(ctx *gin.Context) {
	videoid := ctx.Param("video_id")
	userid := ctx.GetString("authuserid")

	videochan := make(chan *domain.Video, 1)
	errchan := make(chan error, 1)

	go func(videochan chan<- *domain.Video, errchan chan<- error, userid *string, videoid *string) {
		del_video, err := NVh.service.DeleteVideoService(*userid, *videoid)
		if err != nil {
			errchan <- err
			return
		}
		videochan <- del_video
	}(videochan, errchan, &userid, &videoid)
	for {
		select {
		case video := <-videochan:
			ctx.JSON(http.StatusOK, gin.H{
				"success": true,
				"video":   video,
			})
			return
		case err := <-errchan:
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}
	}
}

func (NVh *VideoHandlerStruct) GetVideoDetailsHandler(ctx *gin.Context) {
	videoid := ctx.Param("videoid")

	videochan := make(chan *domain.Video, 1)
	errchan := make(chan error, 1)

	go func(videochan chan<- *domain.Video, errchan chan<- error, videoid *string) {
		video, err := NVh.service.GetVideoDetailsService(*videoid)
		if err != nil {
			errchan <- err
			return
		}
		videochan <- video
	}(videochan, errchan, &videoid)

	for {
		select {
		case video := <-videochan:
			ctx.JSON(
				http.StatusOK,
				gin.H{
					"success": true,
					"data":    video,
				},
			)
			return
		case err := <-errchan:
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

// random videos
func (NVh *VideoHandlerStruct) GetRandomVideosHandler(ctx *gin.Context) {
	limit := ctx.Query("limit")

	videoschan := make(chan *[]domain.Video, 1)
	errchan := make(chan error, 1)

	go func(videoschan chan<- *[]domain.Video, errchan chan<- error, limit *string) {
		videos, err := NVh.service.GetRandomVideosService(*limit)
		if err != nil {
			errchan <- err
			return
		}
		videoschan <- videos
	}(videoschan, errchan, &limit)
	for {
		select {
		case videos := <-videoschan:
			ctx.JSON(
				http.StatusOK,
				gin.H{
					"success": true,
					"data":    videos,
				},
			)
			return
		case err := <-errchan:
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

// search
func (NVh *VideoHandlerStruct) SearchVideoHandler(ctx *gin.Context) {
	query := ctx.Query("query")

	videoschan := make(chan *[]domain.Video, 1)
	errchan := make(chan error, 1)

	go func(videoschan chan<- *[]domain.Video, errchan chan<- error, query *string) {
		videos, err := NVh.service.SearchVideoService(*query)
		if err != nil {
			errchan <- err
			return
		}
		videoschan <- videos
	}(videoschan, errchan, &query)

	for {
		select {
		case videos := <-videoschan:
			ctx.JSON(
				http.StatusOK,
				gin.H{
					"success": true,
					"data":    videos,
				},
			)
			return
		case err := <-errchan:
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

// trending videos
func (NVh *VideoHandlerStruct) TrendingVideosHandler(ctx *gin.Context) {
	videoschan := make(chan *[]domain.Video, 1)
	errorchan := make(chan error, 1)

	go func(videoschan chan<- *[]domain.Video, errchan chan<- error) {
		videos, err := NVh.service.GetTrendingVideoService()
		if err != nil {
			errchan <- err
			return
		}
		videoschan <- videos
	}(videoschan, errorchan)

	for {
		select {
		case videos := <-videoschan:
			ctx.JSON(
				http.StatusOK,
				gin.H{
					"success": true,
					"data":    videos,
				},
			)
			return
		case err := <-errorchan:
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

func (NVh *VideoHandlerStruct) LikeVideoHandler(ctx *gin.Context) {
	userid := ctx.GetString("authuserid")
	videoid := ctx.Param("videoid")

	videochan := make(chan *domain.Video, 1)
	errchan := make(chan error, 1)

	go func(videochan chan<- *domain.Video, errchan chan<- error, userid string, videoid string) {
		liked_video, err := NVh.service.LikeVideoService(userid, videoid)
		if err != nil {
			errchan <- err
			return
		}
		videochan <- liked_video
	}(videochan, errchan, userid, videoid)
	for {
		select {
		case video := <-videochan:
			ctx.JSON(
				http.StatusOK,
				gin.H{
					"success": true,
					"data":    video,
				},
			)
			return
		case err := <-errchan:
			ctx.JSON(
				http.StatusBadRequest,
				gin.H{
					"success": false,
					"error":   err.Error(),
				},
			)
		}
	}
}

func (NVh *VideoHandlerStruct) DislikeVideoHandler(ctx *gin.Context) {
	userid := ctx.GetString("authuserid")
	videoid := ctx.Param("videoid")

	videochan := make(chan *domain.Video, 1)
	errchan := make(chan error, 1)

	go func(videochan chan<- *domain.Video, errchan chan<- error, userid *string, videoid *string) {
		disliked_video, err := NVh.service.DislikeVideoService(*userid, *videoid)
		if err != nil {
			errchan <- err
			return
		}
		videochan <- disliked_video
	}(videochan, errchan, &userid, &videoid)
	for {
		select {
		case video := <-videochan:
			ctx.JSON(
				http.StatusOK,
				gin.H{
					"success": true,
					"data":    video,
				},
			)
			return
		case err := <-errchan:
			ctx.JSON(
				http.StatusBadRequest,
				gin.H{
					"success": false,
					"error":   err.Error(),
				},
			)
			return
		}
	}
}
