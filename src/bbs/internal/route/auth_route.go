package route

import (
	"bbs/internal/controller"
	"bbs/internal/repository"
	"bbs/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetAuthRoute(r *gin.Engine, db *gorm.DB) {
	authRouter := r.Group("/auth")

	authRepository := repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepository)
	authController := controller.NewAuthContorller(authService)

	authRouter.POST("/signup", authController.Signup)
	authRouter.POST("/login", authController.Login)
}
