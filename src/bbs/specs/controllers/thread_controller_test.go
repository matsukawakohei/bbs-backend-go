package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"

	"bbs/controllers"
	"bbs/infra"
	"bbs/middlewares"
	"bbs/models"
	"bbs/repositories"
	"bbs/services"
)

var r *gin.Engine

var db *gorm.DB

var user models.User

type ListResponseBody struct {
	Data []models.Thread `json:"data"`
}

func TestThread(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Thread Suite")
}

var _ = BeforeSuite(func() {
	/** TODO: テスト用の環境を読み込むように修正する */
	infra.Init()
	db = infra.SetUpDB()
	/** ここまで **/

	/** TODO: Routingの設定は関数として切り出す */
	r = gin.Default()
	authRepository := repositories.NewAuthRepository(db)
	authService := services.NewAuthService(authRepository)
	threadRepository := repositories.NewThreadRepository(db)
	threadService := services.NewThreadService(threadRepository)
	threadController := controllers.NewThreadController(threadService)
	threadRouter := r.Group("/threads")
	threadRouterWithAuth := r.Group("/threads", middlewares.AuthMiddleware(authService))

	threadRouter.GET("", threadController.FindAll)
	threadRouter.GET("/:threadId", threadController.FindById)
	threadRouterWithAuth.POST("", threadController.Create)
	threadRouterWithAuth.PUT("/:threadId", threadController.Update)
	threadRouterWithAuth.DELETE("/:threadId", threadController.Delete)
	/** ここまで **/

	/** TODO: ユーザーの作成は関数として切り出す */
	testUserName := "test"
	testUserEmail := "exmaple@example.com"
	testUserPassword := "password"
	user = models.User{
		Name:     testUserName,
		Email:    testUserEmail,
		Password: testUserPassword,
	}
	db.Create(&user)
	/** ここまで **/
})

var _ = Describe("ThreadController", func() {

	Describe("スレッド一覧表示", func() {
		It("スレッドがない場合は空配列を返す", func() {
			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, "/threads", nil)
			r.ServeHTTP(w, req)

			var body ListResponseBody
			decodeErr := json.Unmarshal(w.Body.Bytes(), &body)

			Expect(err).To(BeNil())
			Expect(decodeErr).To(BeNil())
			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(len(body.Data)).To(Equal(0))
		})

		It("スレッドがある場合はスレッドのスライスを返す", func() {
			/** TODO: テストデータの作成は関数として切り出す */
			testTitle := "テストタイトル"
			testBody := "テスト本文"
			testTitleNum := 3

			for i := 0; i < testTitleNum; i++ {
				thread := models.Thread{
					UserID: user.ID,
					Title:  testTitle + strconv.Itoa(i+1),
					Body:   testBody + strconv.Itoa(i+1),
				}
				db.Create(&thread)
			}
			/** ここまで **/

			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, "/threads", nil)
			r.ServeHTTP(w, req)

			var body ListResponseBody
			decoder := json.NewDecoder(bytes.NewReader(w.Body.Bytes()))
			decoder.Decode(&body)

			Expect(err).To(BeNil())
			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(len(body.Data)).To(Equal(3))
		})
	})
})

var _ = AfterSuite(func() {
	/** TODO: テストスイートごとにトランザクションをはれないか調べる */
	db.Where("id > ?", 0).Unscoped().Delete(&models.Thread{})
	db.Unscoped().Delete(&user)
	/** ここまで **/
})
