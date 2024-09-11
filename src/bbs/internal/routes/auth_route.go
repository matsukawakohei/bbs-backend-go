package routes

import (
	"bbs/internal/controller"
	"bbs/internal/repositories"
	"bbs/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetAuthRoute(r *gin.Engine, db *gorm.DB) {
	authRouter := r.Group("/auth")

	authRepository := repositories.NewAuthRepository(db)
	authService := services.NewAuthService(authRepository)
	authController := controller.NewAuthContorller(authService)

	authRouter.POST("/signup", authController.Signup)
	authRouter.POST("/login", authController.Login)
}
