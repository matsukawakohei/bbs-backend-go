package main

import (
	"bbs/infra"
	"bbs/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	infra.Init()
	db := infra.SetUpDB()

	r := gin.Default()

	routes.SetThreadRoute(r, db)
	routes.SetAuthRoute(r, db)
	routes.SetCommentRoute(r, db)

	r.Run("0.0.0.0:8888")
}
