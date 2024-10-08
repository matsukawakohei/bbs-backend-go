package route

import (
	"bbs/internal/controller"
	"bbs/internal/middleware"
	"bbs/internal/repository"
	"bbs/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetThreadRoute(r *gin.Engine, db *gorm.DB) {
	threadRouter := r.Group("/threads")

	authRepository := repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepository)

	threadRouterWithAuth := r.Group("/threads", middleware.AuthMiddleware(authService))

	threadRepository := repository.NewThreadRepository(db)
	threadService := service.NewThreadService(threadRepository)
	threadController := controller.NewThreadController(threadService)

	threadRouter.GET("", threadController.FindAll)
	threadRouter.GET("/:threadId", threadController.FindById)
	threadRouterWithAuth.POST("", threadController.Create)
	threadRouterWithAuth.PUT("/:threadId", threadController.Update)
	threadRouterWithAuth.DELETE("/:threadId", threadController.Delete)
}
