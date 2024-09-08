package routes

import (
	"bbs/internal/controllers"
	"bbs/middlewares"
	"bbs/repositories"
	"bbs/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetCommentRoute(r *gin.Engine, db *gorm.DB) {
	authRepository := repositories.NewAuthRepository(db)
	authService := services.NewAuthService(authRepository)

	threadRepository := repositories.NewThreadRepository(db)

	commentRepository := repositories.NewCommentRepository(db)
	commentService := services.NewCommentService(commentRepository, threadRepository)
	commentController := controllers.NewCommentController(commentService)

	commentRouterWithAuth := r.Group("/threads/:threadId/comments", middlewares.AuthMiddleware(authService))

	commentRouterWithAuth.GET("", commentController.FindByThreadId)
	commentRouterWithAuth.GET("/:commentId", commentController.FindById)
	commentRouterWithAuth.POST("", commentController.Create)
	commentRouterWithAuth.PUT("/:commentId", commentController.Update)
}
