package controller

import (
	"bbs/internal/dto"
	"bbs/internal/models"
	"bbs/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type IThreadController interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	FindAll(ctx *gin.Context)
	FindById(ctx *gin.Context)
}

type ThreadController struct {
	service services.IThreadService
}

func NewThreadController(service services.IThreadService) IThreadController {
	return &ThreadController{service: service}
}

func (c *ThreadController) Create(ctx *gin.Context) {
	user, exists := ctx.Get("user")
	if !exists {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	userId := user.(*models.User).ID

	var input dto.CreateThreadInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newThread, err := c.service.Create(input, userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": newThread})
}

func (c *ThreadController) Update(ctx *gin.Context) {
	user, exists := ctx.Get("user")
	if !exists {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	userId := user.(*models.User).ID

	threadId, err := strconv.ParseUint(ctx.Param("threadId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	var input dto.UpdateThreadInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateThread, err := c.service.Update(uint(threadId), input, userId)
	if err != nil {
		if err.Error() == "user is not thread owner" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if err.Error() == "thread not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": updateThread})
}

func (c *ThreadController) Delete(ctx *gin.Context) {
	user, exists := ctx.Get("user")
	if !exists {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	userId := user.(*models.User).ID

	threadId, err := strconv.ParseUint(ctx.Param("threadId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	err = c.service.Delete(uint(threadId), userId)
	if err != nil {
		if err.Error() == "user is not thread owner" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if err.Error() == "thread not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected error"})
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *ThreadController) FindAll(ctx *gin.Context) {
	threads, err := c.service.FindAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": threads})
}

func (c *ThreadController) FindById(ctx *gin.Context) {
	threadId, err := strconv.ParseUint(ctx.Param("threadId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	thread, err := c.service.FindById(uint(threadId))
	if err != nil {
		if err.Error() == "thread not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "unexpected error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": thread})
}
