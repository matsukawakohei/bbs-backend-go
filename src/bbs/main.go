package main

import (
	"bbs/controllers"
	"bbs/infra"
	"bbs/middlewares"
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

	threadRepository := repositories.NewThreadRepository(db)
	threadService := services.NewThreadService(threadRepository)
	threadController := controllers.NewThreadController(threadService)

	r := gin.Default()
	authRouter := r.Group("/auth")
	threadRouter := r.Group("/threads")
	threadRouterWithAuth := r.Group("/threads", middlewares.AuthMiddleware(authService))

	r.GET("/sample", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	authRouter.POST("/signup", authController.Signup)
	authRouter.POST("/login", authController.Login)

	threadRouter.GET("", threadController.FindAll)
	threadRouter.GET("/:id", threadController.FindById)
	threadRouterWithAuth.POST("", threadController.Create)
	threadRouterWithAuth.PUT("/:id", threadController.Update)
	threadRouterWithAuth.DELETE("/:id", threadController.Delete)
	r.Run("0.0.0.0:8888")
}
