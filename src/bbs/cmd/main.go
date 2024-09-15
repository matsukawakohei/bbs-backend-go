package main

import (
	"bbs/internal/infra"
	"bbs/internal/route"

	"github.com/gin-gonic/gin"
)

func main() {
	infra.Init()
	db := infra.SetUpDB()

	r := gin.Default()

	route.SetCorsHeader(r)

	route.SetThreadRoute(r, db)
	route.SetAuthRoute(r, db)
	route.SetCommentRoute(r, db)

	r.Run("0.0.0.0:8888")
}
