package route

import (
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetCorsHeader(r *gin.Engine) {
	r.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Split(os.Getenv("FRONT_URL"), ","), // Nuxt.jsのオリジン
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},   // 許可するHTTPメソッド
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true, // クッキーや認証情報を許可する場合
		MaxAge:           12 * time.Hour,
	}))
}
