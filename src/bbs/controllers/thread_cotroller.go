package controllers

import (
	"bbs/dto"
	"bbs/models"
	"bbs/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IThreadController interface {
	Create(ctx *gin.Context)
	FindAll(ctx *gin.Context)
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

func (c *ThreadController) FindAll(ctx *gin.Context) {
	threads, err := c.service.FindAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": threads})
}
