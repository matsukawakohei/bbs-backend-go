package routes

import (
	"bbs/internal/controllers"
	"bbs/internal/middlewares"
	"bbs/internal/repositories"
	"bbs/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetThreadRoute(r *gin.Engine, db *gorm.DB) {
	threadRouter := r.Group("/threads")

	authRepository := repositories.NewAuthRepository(db)
	authService := services.NewAuthService(authRepository)

	threadRouterWithAuth := r.Group("/threads", middlewares.AuthMiddleware(authService))

	threadRepository := repositories.NewThreadRepository(db)
	threadService := services.NewThreadService(threadRepository)
	threadController := controllers.NewThreadController(threadService)

	threadRouter.GET("", threadController.FindAll)
	threadRouter.GET("/:threadId", threadController.FindById)
	threadRouterWithAuth.POST("", threadController.Create)
	threadRouterWithAuth.PUT("/:threadId", threadController.Update)
	threadRouterWithAuth.DELETE("/:threadId", threadController.Delete)
}
