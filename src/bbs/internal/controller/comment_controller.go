package controller

import (
	"bbs/internal/dto"
	"bbs/internal/model"
	"bbs/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommentController struct {
	service service.ICommentService
}

func NewCommentController(service service.ICommentService) ICommentController {
	return &CommentController{service: service}
}

func (c *CommentController) Create(ctx *gin.Context) {
	user, exists := ctx.Get("user")
	if !exists {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	userId := user.(*model.User).ID

	threadId, err := strconv.ParseUint(ctx.Param("threadId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	var input dto.CreateComment
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newComment, err := c.service.Create(input, uint(threadId), userId)
	if err != nil {
		if err.Error() == "thread not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected error"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": newComment})
}

func (c *CommentController) Update(ctx *gin.Context) {
	user, exists := ctx.Get("user")
	if !exists {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	userId := user.(*model.User).ID

	threadId, err := strconv.ParseUint(ctx.Param("threadId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid thread id"})
		return
	}

	commentId, err := strconv.ParseUint(ctx.Param("commentId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment id"})
		return
	}

	var input dto.UpdateComment
	if err = ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateComment, err := c.service.Update(input, uint(commentId), uint(threadId), userId)
	if err != nil {
		if err.Error() == "user is not comment owner" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		}
		if err.Error() == "comment not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": updateComment})
}

func (c *CommentController) Delete(ctx *gin.Context) {
	user, exists := ctx.Get("user")
	if !exists {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	userId := user.(*model.User).ID

	threadId, err := strconv.ParseUint(ctx.Param("threadId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid thread id"})
		return
	}

	commentId, err := strconv.ParseUint(ctx.Param("commentId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment id"})
		return
	}

	err = c.service.Delete(uint(commentId), uint(threadId), userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected error"})
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *CommentController) FindByThreadId(ctx *gin.Context) {
	user, exists := ctx.Get("user")
	if !exists {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	userId := user.(*model.User).ID

	threadId, err := strconv.ParseUint(ctx.Param("threadId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid thread id"})
		return
	}

	comments, err := c.service.FindByThreadId(uint(threadId), userId)
	if err != nil {
		if err.Error() == "comment not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": comments})
}

func (c *CommentController) FindById(ctx *gin.Context) {
	threadId, err := strconv.ParseUint(ctx.Param("threadId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid thread id"})
		return
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
	}

	comment, err := c.service.FindById(uint(id), uint(threadId))
	if err != nil {
		if err.Error() == "comment not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected error"})
	}

	ctx.JSON(http.StatusOK, gin.H{"data": comment})
}
