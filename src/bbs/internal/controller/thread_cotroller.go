package controller

import (
	"bbs/internal/dto"
	"bbs/internal/model"
	"bbs/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ThreadController struct {
	service service.IThreadService
}

func NewThreadController(service service.IThreadService) IThreadController {
	return &ThreadController{service: service}
}

func (c *ThreadController) Create(ctx *gin.Context) {
	user, exists := ctx.Get("user")
	if !exists {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	userId := user.(*model.User).ID

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

	userId := user.(*model.User).ID

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

	userId := user.(*model.User).ID

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
	// TODO: デフォルトのlimitは別ファイルで定数として定義する
	limit := 10
	if limitQuery := ctx.Query("limit"); limitQuery != "" {
		limitInt, err := strconv.Atoi(limitQuery)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		limit = limitInt
	}

	// TODO: デフォルトのoffsetは別ファイルで定数として定義する
	offset := 0
	if pageQuery := ctx.Query("page"); pageQuery != "" {
		pageInt, err := strconv.Atoi(pageQuery)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// offsetは0始まりでpageは1始まりなので1を引く
		offset = pageInt - 1
	}

	threadList, err := c.service.FindAll(limit, offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": threadList})
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
