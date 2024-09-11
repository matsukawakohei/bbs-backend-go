package routes

import (
	"bbs/internal/controller"
	"bbs/internal/middleware"
	"bbs/internal/repositories"
	"bbs/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetCommentRoute(r *gin.Engine, db *gorm.DB) {
	authRepository := repositories.NewAuthRepository(db)
	authService := services.NewAuthService(authRepository)

	threadRepository := repositories.NewThreadRepository(db)

	commentRepository := repositories.NewCommentRepository(db)
	commentService := services.NewCommentService(commentRepository, threadRepository)
	commentController := controller.NewCommentController(commentService)

	commentRouterWithAuth := r.Group("/threads/:threadId/comments", middleware.AuthMiddleware(authService))

	commentRouterWithAuth.GET("", commentController.FindByThreadId)
	commentRouterWithAuth.GET("/:commentId", commentController.FindById)
	commentRouterWithAuth.POST("", commentController.Create)
	commentRouterWithAuth.PUT("/:commentId", commentController.Update)
}
