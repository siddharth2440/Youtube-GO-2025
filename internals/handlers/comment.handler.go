package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/itsmonday/youtube/internals/domain"
	"github.com/itsmonday/youtube/internals/services"
)

type CommentHandlerStruct struct {
	services services.CommentServiceInterface
}

func NewCommentHandler(services services.CommentServiceInterface) *CommentHandlerStruct {
	return &CommentHandlerStruct{services}
}

// add a comment
func (NCh *CommentHandlerStruct) AddCommentHandler(ctx *gin.Context) {
	var comment domain.Comment
	if err := ctx.ShouldBindJSON(&comment); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"success": false,
				"error":   err.Error(),
			},
		)
		return
	}

	videoId := ctx.Param("videoid")
	userId := ctx.GetString("authuserid")

	commentchan := make(chan *domain.Comment, 1)
	errchan := make(chan error, 1)

	go func(commentchan chan<- *domain.Comment, errchan chan<- error, videoid string, userid string) {
		newcomment, err := NCh.services.AddCommentService(&comment, userid, videoid)
		if err != nil {
			errchan <- err
			return
		}
		commentchan <- newcomment
	}(commentchan, errchan, videoId, userId)

	for {
		select {
		case comment := <-commentchan:
			ctx.JSON(
				http.StatusOK,
				gin.H{
					"success": true,
					"data":    comment,
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

// get comments of a video
func (NCh *CommentHandlerStruct) GetCommentsHandler(ctx *gin.Context) {
	videoid := ctx.Param("videoid")

	commentschan := make(chan *[]domain.Comment, 1)
	errchan := make(chan error, 1)

	go func(commentschan chan<- *[]domain.Comment, errchan chan<- error, videoid string) {
		comments, err := NCh.services.GetCommentsService(videoid)
		if err != nil {
			errchan <- err
			return
		}
		commentschan <- comments
	}(commentschan, errchan, videoid)

	for {
		select {
		case comments := <-commentschan:
			ctx.JSON(
				http.StatusOK,
				gin.H{
					"success": true,
					"data":    comments,
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

// get a particular comment
func (NCh *CommentHandlerStruct) GetCommentDetailsHandler(ctx *gin.Context) {
	commentid := ctx.Param("commentid")

	commentchan := make(chan *domain.Comment, 1)
	errchan := make(chan error, 1)

	go func(commentchan chan<- *domain.Comment, errchan chan<- error, commentid string) {
		comment, err := NCh.services.GetcommentDetails(commentid)
		if err != nil {
			errchan <- err
			return
		}
		commentchan <- comment
	}(commentchan, errchan, commentid)

	for {
		select {
		case comment := <-commentchan:
			ctx.JSON(
				http.StatusAccepted,
				gin.H{
					"success": true,
					"data":    comment,
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

// delete comment
func (NCh *CommentHandlerStruct) DeleteCommentHandler(ctx *gin.Context) {
	commentid := ctx.Param("commentid")
	userid := ctx.GetString("authuserid")

	commentchan := make(chan *domain.Comment, 1)
	errchan := make(chan error, 1)

	go func(commentchan chan<- *domain.Comment, errchan chan<- error, commentid, userid string) {
		comment, err := NCh.services.DeleteCommentService(commentid, userid)
		if err != nil {
			errchan <- err
			return
		}
		commentchan <- comment
	}(commentchan, errchan, commentid, userid)

	for {
		select {
		case comment := <-commentchan:
			ctx.JSON(
				http.StatusOK,
				gin.H{
					"success": true,
					"data":    comment,
					"mesage":  "Comment deleted succesfully",
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
