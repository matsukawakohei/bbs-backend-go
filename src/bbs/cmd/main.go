package main

import (
	"bbs/internal/infra"
	"bbs/internal/route"
	"os"

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

	allowHost := os.Getenv("ALLOW_HOST")
	port := os.Getenv("PORT")

	r.Run(allowHost + ":" + port)
}
