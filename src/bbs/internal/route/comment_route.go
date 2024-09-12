package route

import (
	"bbs/internal/controller"
	"bbs/internal/middleware"
	"bbs/internal/repository"
	"bbs/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetCommentRoute(r *gin.Engine, db *gorm.DB) {
	authRepository := repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepository)

	threadRepository := repository.NewThreadRepository(db)

	commentRepository := repository.NewCommentRepository(db)
	commentService := service.NewCommentService(commentRepository, threadRepository)
	commentController := controller.NewCommentController(commentService)

	commentRouterWithAuth := r.Group("/threads/:threadId/comments", middleware.AuthMiddleware(authService))

	commentRouterWithAuth.GET("", commentController.FindByThreadId)
	commentRouterWithAuth.GET("/:commentId", commentController.FindById)
	commentRouterWithAuth.POST("", commentController.Create)
	commentRouterWithAuth.PUT("/:commentId", commentController.Update)
}
