package controller_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"

	"bbs/internal/dto"
	"bbs/internal/model"
)

type ListResponseBody struct {
	Data dto.ThreadListOutput `json:"data"`
}

type CreateRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type CreateResponse struct {
	Thread       model.Thread `json:"data"`
	ErrorMessage string       `json:"error"`
}

type DetailResponse struct {
	Thread model.Thread `json:"data"`
}

type UpdateRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type UpdateResponse struct {
	Thread       model.Thread `json:"data"`
	ErrorMessage string       `json:"error"`
}

var _ = Describe("ThreadController", func() {
	BeforeEach(func() {
		defaultBeforeEachFunc()
	})

	AfterEach(func() {
		defaultAfterEachFunc()
	})

	Describe("スレッド一覧表示", func() {
		Context("スレッドがない場合", func() {
			It("空配列を返す", func() {
				w := httptest.NewRecorder()
				req, err := http.NewRequest(http.MethodGet, "/threads", nil)
				r.ServeHTTP(w, req)

				var body ListResponseBody
				decodeErr := json.Unmarshal(w.Body.Bytes(), &body)

				Expect(err).To(BeNil())
				Expect(decodeErr).To(BeNil())
				Expect(w.Code).To(Equal(http.StatusOK))
				fmt.Println(body.Data)
				Expect(body.Data.Total).To(Equal(int64(0)))
				Expect(len(body.Data.Threads)).To(Equal(0))
			})
		})

		Context("スレッドがある場合", func() {
			It("1ページ分スレッドのスライスと合計件数を返す", func() {
				testThreadNum := 10
				createTestThread(db, user.ID, testThreadNum)

				w := httptest.NewRecorder()
				req, err := http.NewRequest(http.MethodGet, "/threads?page=1&limit=5", nil)
				r.ServeHTTP(w, req)

				var body ListResponseBody
				decoder := json.NewDecoder(bytes.NewReader(w.Body.Bytes()))
				decoder.Decode(&body)

				Expect(err).To(BeNil())
				Expect(w.Code).To(Equal(http.StatusOK))
				Expect(body.Data.Total).To(Equal(int64(10)))
				Expect(len(body.Data.Threads)).To(Equal(5))
			})
		})

		Context("ページネーションの指定", func() {
			It("指定ページ分スレッドのスライスと合計件数を返す(ページ指定)", func() {
				testThreadNum := 6
				createTestThread(db, user.ID, testThreadNum)

				w := httptest.NewRecorder()
				req, err := http.NewRequest(http.MethodGet, "/threads?page=2&limit=5", nil)
				r.ServeHTTP(w, req)

				var body ListResponseBody
				decoder := json.NewDecoder(bytes.NewReader(w.Body.Bytes()))
				decoder.Decode(&body)

				Expect(err).To(BeNil())
				Expect(w.Code).To(Equal(http.StatusOK))
				Expect(body.Data.Total).To(Equal(int64(6)))
				Expect(len(body.Data.Threads)).To(Equal(1))
			})
		})
	})

	Describe("スレッド詳細取得", func() {
		It("スレッド詳細を取得する", func() {
			testThreadNum := 1
			testThread := createTestThread(db, user.ID, testThreadNum)[0]

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
			db.Where("id > ?", 0).Unscoped().Delete(&model.Comment{})
			db.Where("id > ?", 0).Unscoped().Delete(&model.Thread{})
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
			testThread := createTestThread(db, user.ID, testThreadNum)[0]

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
			testThread := createTestThread(db, user.ID, testThreadNum)[0]

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
			testThread := createTestThread(db, user.ID, testThreadNum)[0]

			name := "test"
			email := "exampleexample@example.com"
			otherUser := createTestUser(r, db, name, email)
			otherUserToken := createTestUserToken(r, otherUser.Email)

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
			testThread := createTestThread(db, user.ID, testThreadNum)[0]

			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodDelete, "/threads/"+strconv.Itoa(int(testThread.ID)), nil)
			req.Header.Set("Content-Type", contentType)
			req.Header.Set("Authorization", "Bearer "+token)
			r.ServeHTTP(w, req)

			Expect(err).To(BeNil())
			Expect(w.Code).To(Equal(http.StatusOK))

			var deletedThread model.Thread
			result := db.First(&deletedThread, testThread.ID)
			Expect(errors.Is(result.Error, gorm.ErrRecordNotFound)).To(BeTrue())
		})

		It("トークンがなければエラー", func() {
			testThreadNum := 1
			testThread := createTestThread(db, user.ID, testThreadNum)[0]

			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodDelete, "/threads/"+strconv.Itoa(int(testThread.ID)), nil)
			req.Header.Set("Content-Type", contentType)
			r.ServeHTTP(w, req)

			Expect(err).To(BeNil())
			Expect(w.Code).To(Equal(http.StatusUnauthorized))
		})

		It("スレッドの所有者ではない場合は削除できない", func() {
			testThreadNum := 1
			testThread := createTestThread(db, user.ID, testThreadNum)[0]

			name := "test"
			email := "exampleexample@example.com"
			otherUser := createTestUser(r, db, name, email)
			otherUserToken := createTestUserToken(r, otherUser.Email)

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
