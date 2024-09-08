package controllers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"

	"bbs/internal/infra"
	"bbs/internal/routes"
	"bbs/models"
	"bbs/specs/utils"
)

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

type DetailResponse struct {
	Thread models.Thread `json:"data"`
}

type UpdateRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type UpdateResponse struct {
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
	routes.SetCommentRoute(r, db)

	name := "test"
	email := "exmaple@example.com"
	user = utils.CreateTestUser(r, db, name, email)
	token = utils.CreateTestUserToken(r, user.Email)
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

	Describe("スレッド詳細取得", func() {
		It("スレッド詳細を取得する", func() {
			testThreadNum := 1
			testThread := utils.CreateTestThread(db, user.ID, testThreadNum)[0]

			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, "/threads/"+strconv.Itoa(int(testThread.ID)), nil)
			r.ServeHTTP(w, req)

			var res DetailResponse
			decoder := json.NewDecoder(bytes.NewReader(w.Body.Bytes()))
			decoder.Decode(&res)

			Expect(err).To(BeNil())
			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(res.Thread.ID).To(Equal(testThread.ID))
			Expect(res.Thread.Title).To(Equal(testThread.Title))
			Expect(res.Thread.Body).To(Equal(testThread.Body))
			Expect(res.Thread.UserID).To(Equal(testThread.UserID))
		})

		It("スレッドが存在しない場合は404", func() {

			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, "/threads/"+strconv.Itoa(0), nil)
			r.ServeHTTP(w, req)

			var res DetailResponse
			decoder := json.NewDecoder(bytes.NewReader(w.Body.Bytes()))
			decoder.Decode(&res)

			Expect(err).To(BeNil())
			Expect(w.Code).To(Equal(http.StatusNotFound))
		})

		It("パラメータが文字列の場合はバリデーションエラー", func() {

			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, "/threads/aaa", nil)
			r.ServeHTTP(w, req)

			var res DetailResponse
			decoder := json.NewDecoder(bytes.NewReader(w.Body.Bytes()))
			decoder.Decode(&res)

			Expect(err).To(BeNil())
			Expect(w.Code).To(Equal(http.StatusBadRequest))
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

	Describe("スレッド更新", func() {
		AfterEach(func() {
			db.Where("id > ?", 0).Unscoped().Delete(&models.Comment{})
			db.Where("id > ?", 0).Unscoped().Delete(&models.Thread{})
		})

		It("スレッドを更新する", func() {
			title := "test"
			body := "testtest"

			request := UpdateRequest{
				Title: title,
				Body:  body,
			}
			requestBytes, _ := json.Marshal(request)

			testThreadNum := 1
			testThread := utils.CreateTestThread(db, user.ID, testThreadNum)[0]

			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPut, "/threads/"+strconv.Itoa(int(testThread.ID)), bytes.NewBuffer(requestBytes))
			req.Header.Set("Content-Type", contentType)
			req.Header.Set("Authorization", "Bearer "+token)
			r.ServeHTTP(w, req)

			var res UpdateResponse
			decoder := json.NewDecoder(bytes.NewReader(w.Body.Bytes()))
			decoder.Decode(&res)

			Expect(err).To(BeNil())
			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(res.Thread.ID).To(Equal(testThread.ID))
			Expect(res.Thread.Title).To(Equal(title))
			Expect(res.Thread.Body).To(Equal(body))
			Expect(res.Thread.UserID).To(Equal(user.ID))
		})

		It("トークンがなければエラー", func() {
			title := "test"
			body := "testtest"

			request := UpdateRequest{
				Title: title,
				Body:  body,
			}
			requestBytes, _ := json.Marshal(request)

			testThreadNum := 1
			testThread := utils.CreateTestThread(db, user.ID, testThreadNum)[0]

			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPut, "/threads/"+strconv.Itoa(int(testThread.ID)), bytes.NewBuffer(requestBytes))
			req.Header.Set("Content-Type", contentType)
			r.ServeHTTP(w, req)

			Expect(err).To(BeNil())
			Expect(w.Code).To(Equal(http.StatusUnauthorized))
		})

		It("スレッドの所有者ではない場合は更新できない", func() {
			title := "test"
			body := "testtest"

			request := UpdateRequest{
				Title: title,
				Body:  body,
			}
			requestBytes, _ := json.Marshal(request)

			testThreadNum := 1
			testThread := utils.CreateTestThread(db, user.ID, testThreadNum)[0]

			name := "test"
			email := "exampleexample@example.com"
			otherUser := utils.CreateTestUser(r, db, name, email)
			otherUserToken := utils.CreateTestUserToken(r, otherUser.Email)

			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPut, "/threads/"+strconv.Itoa(int(testThread.ID)), bytes.NewBuffer(requestBytes))
			req.Header.Set("Content-Type", contentType)
			req.Header.Set("Authorization", "Bearer "+otherUserToken)
			r.ServeHTTP(w, req)

			db.Unscoped().Delete(&otherUser)

			Expect(err).To(BeNil())
			Expect(w.Code).To(Equal(http.StatusUnauthorized))
		})

		It("スレッドが存在しない場合は404", func() {
			title := "test"
			body := "testtest"

			request := UpdateRequest{
				Title: title,
				Body:  body,
			}
			requestBytes, _ := json.Marshal(request)

			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPut, "/threads/"+strconv.Itoa(0), bytes.NewBuffer(requestBytes))
			req.Header.Set("Content-Type", contentType)
			req.Header.Set("Authorization", "Bearer "+token)
			r.ServeHTTP(w, req)

			Expect(err).To(BeNil())
			Expect(w.Code).To(Equal(http.StatusNotFound))
		})

		It("URLパラメータが文字列の場合はバリデーションエラー", func() {
			title := "test"
			body := "testtest"

			request := UpdateRequest{
				Title: title,
				Body:  body,
			}
			requestBytes, _ := json.Marshal(request)

			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPut, "/threads/aaa", bytes.NewBuffer(requestBytes))
			req.Header.Set("Content-Type", contentType)
			req.Header.Set("Authorization", "Bearer "+token)
			r.ServeHTTP(w, req)

			Expect(err).To(BeNil())
			Expect(w.Code).To(Equal(http.StatusBadRequest))
		})
	})

	Describe("スレッド削除", func() {
		It("スレッドを削除する", func() {
			testThreadNum := 1
			testThread := utils.CreateTestThread(db, user.ID, testThreadNum)[0]

			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodDelete, "/threads/"+strconv.Itoa(int(testThread.ID)), nil)
			req.Header.Set("Content-Type", contentType)
			req.Header.Set("Authorization", "Bearer "+token)
			r.ServeHTTP(w, req)

			Expect(err).To(BeNil())
			Expect(w.Code).To(Equal(http.StatusOK))

			var deletedThread models.Thread
			result := db.First(&deletedThread, testThread.ID)
			Expect(errors.Is(result.Error, gorm.ErrRecordNotFound)).To(BeTrue())
		})

		It("トークンがなければエラー", func() {
			testThreadNum := 1
			testThread := utils.CreateTestThread(db, user.ID, testThreadNum)[0]

			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodDelete, "/threads/"+strconv.Itoa(int(testThread.ID)), nil)
			req.Header.Set("Content-Type", contentType)
			r.ServeHTTP(w, req)

			Expect(err).To(BeNil())
			Expect(w.Code).To(Equal(http.StatusUnauthorized))
		})

		It("スレッドの所有者ではない場合は削除できない", func() {
			testThreadNum := 1
			testThread := utils.CreateTestThread(db, user.ID, testThreadNum)[0]

			name := "test"
			email := "exampleexample@example.com"
			otherUser := utils.CreateTestUser(r, db, name, email)
			otherUserToken := utils.CreateTestUserToken(r, otherUser.Email)

			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodDelete, "/threads/"+strconv.Itoa(int(testThread.ID)), nil)
			req.Header.Set("Content-Type", contentType)
			req.Header.Set("Authorization", "Bearer "+otherUserToken)
			r.ServeHTTP(w, req)

			db.Unscoped().Delete(&otherUser)

			Expect(err).To(BeNil())
			Expect(w.Code).To(Equal(http.StatusUnauthorized))
		})

		It("スレッドが存在しない場合は404", func() {
			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodDelete, "/threads/"+strconv.Itoa(0), nil)
			req.Header.Set("Content-Type", contentType)
			req.Header.Set("Authorization", "Bearer "+token)
			r.ServeHTTP(w, req)

			Expect(err).To(BeNil())
			Expect(w.Code).To(Equal(http.StatusNotFound))
		})

		It("URLパラメータが文字列の場合はバリデーションエラー", func() {
			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodDelete, "/threads/aaa", nil)
			req.Header.Set("Content-Type", contentType)
			req.Header.Set("Authorization", "Bearer "+token)
			r.ServeHTTP(w, req)

			Expect(err).To(BeNil())
			Expect(w.Code).To(Equal(http.StatusBadRequest))
		})
	})
})

var _ = AfterSuite(func() {
	db.Where("id > ?", 0).Unscoped().Delete(&models.Comment{})
	db.Where("id > ?", 0).Unscoped().Delete(&models.Thread{})
	db.Unscoped().Delete(&user)
})
