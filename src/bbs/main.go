package main

import (
	"bbs/controllers"
	"bbs/infra"
	"bbs/repositories"
	"bbs/services"

	"github.com/gin-gonic/gin"
)

func main() {
	infra.Init()
	db := infra.SetUpDB()

	authRepository := repositories.NewAuthRepository(db)
	authService := services.NewAuthService(authRepository)
	authController := controllers.NewAuthContorller(authService)

	r := gin.Default()
	authRouter := r.Group("/auth")

	r.GET("/sample", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	authRouter.POST("/signup", authController.Signup)
	authRouter.POST("/login", authController.Login)
	r.Run("0.0.0.0:8888")
}
