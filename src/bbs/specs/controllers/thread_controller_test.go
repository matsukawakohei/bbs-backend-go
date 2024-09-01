package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"

	"bbs/infra"
	"bbs/models"
	"bbs/routes"
	"bbs/specs/utils"
)

var r *gin.Engine

var db *gorm.DB

var user *models.User

var token string

type ListResponseBody struct {
	Data []models.Thread `json:"data"`
}

type CreateRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type CreateResponse struct {
	Thread       models.Thread `json:"data"`
	ErrorMessage string        `json:"error"`
}

func TestThread(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Thread Suite")
}

var _ = BeforeSuite(func() {
	infra.TestInit(utils.GetEnvTestPath())
	db = infra.SetUpDB()

	r = gin.New()
	routes.SetThreadRoute(r, db)
	routes.SetAuthRoute(r, db)

	user = utils.CreateTestUser(r, db)
	token = utils.CreateTestUserToken(r, user.Email)
})

var _ = Describe("ThreadController", func() {

	contentType := "application/json"

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
			testThreadNum := 3
			utils.CreateTestThread(db, user.ID, testThreadNum)

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

	Describe("スレッド作成", func() {
		It("スレッドを作成する", func() {
			title := "テスト"
			body := "テストテスト"

			request := CreateRequest{
				Title: title,
				Body:  body,
			}
			requestBytes, _ := json.Marshal(request)

			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, "/threads", bytes.NewBuffer(requestBytes))
			req.Header.Set("Content-Type", contentType)
			req.Header.Set("Authorization", "Bearer "+token)
			r.ServeHTTP(w, req)

			var res CreateResponse
			decoder := json.NewDecoder(bytes.NewReader(w.Body.Bytes()))
			decoder.Decode(&res)

			Expect(err).To(BeNil())
			Expect(w.Code).To(Equal(http.StatusCreated))
			Expect(res.Thread.ID).NotTo(BeNil())
			Expect(res.Thread.Title).To(Equal(title))
			Expect(res.Thread.Body).To(Equal(body))
			Expect(res.Thread.UserID).To(Equal(user.ID))
		})

		It("トークンがなければエラー", func() {
			title := "テスト"
			body := "テストテスト"

			request := CreateRequest{
				Title: title,
				Body:  body,
			}
			requestBytes, _ := json.Marshal(request)

			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, "/threads", bytes.NewBuffer(requestBytes))
			req.Header.Set("Content-Type", contentType)
			r.ServeHTTP(w, req)

			Expect(err).To(BeNil())
			Expect(w.Code).To(Equal(http.StatusUnauthorized))
		})

		It("タイトルがない場合はエラー", func() {
			title := "テスト"

			request := CreateRequest{
				Title: title,
			}
			requestBytes, _ := json.Marshal(request)
			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, "/threads", bytes.NewBuffer(requestBytes))
			req.Header.Set("Content-Type", contentType)
			req.Header.Set("Authorization", "Bearer "+token)
			r.ServeHTTP(w, req)

			var res CreateResponse
			decoder := json.NewDecoder(bytes.NewReader(w.Body.Bytes()))
			decoder.Decode(&res)

			Expect(err).To(BeNil())
			Expect(w.Code).To(Equal(http.StatusBadRequest))
			Expect(res.ErrorMessage).To(ContainSubstring("failed on the 'required' tag"))
		})

		It("本文がない場合はエラー", func() {
			body := "テストテスト"

			request := CreateRequest{
				Body: body,
			}
			requestBytes, _ := json.Marshal(request)
			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, "/threads", bytes.NewBuffer(requestBytes))
			req.Header.Set("Content-Type", contentType)
			req.Header.Set("Authorization", "Bearer "+token)
			r.ServeHTTP(w, req)

			var res CreateResponse
			decoder := json.NewDecoder(bytes.NewReader(w.Body.Bytes()))
			decoder.Decode(&res)

			Expect(err).To(BeNil())
			Expect(w.Code).To(Equal(http.StatusBadRequest))
			Expect(res.ErrorMessage).To(ContainSubstring("failed on the 'required' tag"))
		})
	})
})

var _ = AfterSuite(func() {
	db.Where("id > ?", 0).Unscoped().Delete(&models.Thread{})
	db.Unscoped().Delete(&user)
})
