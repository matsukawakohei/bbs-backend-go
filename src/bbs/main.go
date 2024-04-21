package main

import (
	"bbs/infra"

	"github.com/gin-gonic/gin"
)

func main() {
	infra.Init()
	r := gin.Default()
	r.GET("/sample", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run("0.0.0.0:8888")
}
