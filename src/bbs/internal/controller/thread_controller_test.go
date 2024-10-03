package controller_test

import (
	"bytes"
	"encoding/json"
	"errors"
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
			It("ステータスコード200を返す", func() {
				url := "/threads"
				w := requestAPI(http.MethodGet, url, "", nil)

				Expect(w.Code).To(Equal(http.StatusOK))
			})

			It("空配列を返す", func() {
				url := "/threads"
				w := requestAPI(http.MethodGet, url, "", nil)

				body := getThreadListResponseBody(w)

				Expect(body.Data.Total).To(Equal(int64(0)))
				Expect(len(body.Data.Threads)).To(Equal(0))
			})
		})

		Context("スレッドがある場合", func() {
			It("ステータスコード200を返す", func() {
				testThreadNum := 10
				createTestThread(db, user.ID, testThreadNum)

				url := "/threads?page=1&limit=5"
				w := requestAPI(http.MethodGet, url, "", nil)

				Expect(w.Code).To(Equal(http.StatusOK))
			})

			It("1ページ分スレッドのスライスと合計件数を返す", func() {
				testThreadNum := 10
				createTestThread(db, user.ID, testThreadNum)

				url := "/threads?page=1&limit=5"
				w := requestAPI(http.MethodGet, url, "", nil)

				body := getThreadListResponseBody(w)

				Expect(body.Data.Total).To(Equal(int64(10)))
				Expect(len(body.Data.Threads)).To(Equal(5))
			})
		})

		Context("ページネーションの指定", func() {
			It("ステータスコード200を返す", func() {
				testThreadNum := 6
				createTestThread(db, user.ID, testThreadNum)

				url := "/threads?page=2&limit=5"
				w := requestAPI(http.MethodGet, url, "", nil)

				Expect(w.Code).To(Equal(http.StatusOK))
			})

			It("指定ページ分スレッドのスライスと合計件数を返す(ページ指定)", func() {
				testThreadNum := 6
				createTestThread(db, user.ID, testThreadNum)

				url := "/threads?page=2&limit=5"
				w := requestAPI(http.MethodGet, url, "", nil)

				body := getThreadListResponseBody(w)

				Expect(body.Data.Total).To(Equal(int64(6)))
				Expect(len(body.Data.Threads)).To(Equal(1))
			})
		})
	})

	Describe("スレッド詳細取得", func() {
		Context("スレッドが存在する場合", func() {
			It("ステータスコード200を返す", func() {
				testThreadNum := 1
				testThread := createTestThread(db, user.ID, testThreadNum)[0]

				url := "/threads/" + strconv.Itoa(int(testThread.ID))
				w := requestAPI(http.MethodGet, url, "", nil)

				Expect(w.Code).To(Equal(http.StatusOK))
			})

			It("スレッドを返す", func() {
				testThreadNum := 1
				testThread := createTestThread(db, user.ID, testThreadNum)[0]

				url := "/threads/" + strconv.Itoa(int(testThread.ID))
				w := requestAPI(http.MethodGet, url, "", nil)

				body := getThreadDetailResponseBody(w)

				Expect(body.Thread.ID).To(Equal(testThread.ID))
				Expect(body.Thread.Title).To(Equal(testThread.Title))
				Expect(body.Thread.Body).To(Equal(testThread.Body))
				Expect(body.Thread.UserID).To(Equal(testThread.UserID))
			})
		})

		Context("スレッドが存在しない場合", func() {
			It("エラーコード404を返す", func() {
				url := "/threads/" + strconv.Itoa(0)
				w := requestAPI(http.MethodGet, url, "", nil)

				Expect(w.Code).To(Equal(http.StatusNotFound))
			})
		})

		Context("パラメータが文字列の場合", func() {
			It("エラーコード400を返す", func() {
				url := "/threads/aaa"
				w := requestAPI(http.MethodGet, url, "", nil)

				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
		})
	})

	Describe("スレッド作成", func() {
		Context("スレッドを作成できた場合", func() {
			It("ステータスコード201を返す", func() {
				title := "テスト"
				body := "テストテスト"

				request := getCreateThreadRequestBodyBites(title, body)

				url := "/threads"
				w := requestAPI(http.MethodPost, url, token, request)

				Expect(w.Code).To(Equal(http.StatusCreated))
			})

			It("作成したスレッドを返す", func() {
				title := "テスト"
				body := "テストテスト"

				request := getCreateThreadRequestBodyBites(title, body)

				url := "/threads"
				w := requestAPI(http.MethodPost, url, token, request)

				responseBody := getThreadCreateResponseBody(w)

				Expect(responseBody.Thread.ID).NotTo(BeNil())
				Expect(responseBody.Thread.Title).To(Equal(title))
				Expect(responseBody.Thread.Body).To(Equal(body))
				Expect(responseBody.Thread.UserID).To(Equal(user.ID))
			})

			It("作成したスレッドがDBに保存されている", func() {
				title := "テスト"
				body := "テストテスト"

				request := getCreateThreadRequestBodyBites(title, body)

				url := "/threads"
				requestAPI(http.MethodPost, url, token, request)

				var dbThread model.Thread
				result := db.First(&dbThread)

				Expect(result.Error).To(BeNil())
				Expect(dbThread.ID).NotTo(BeNil())
				Expect(dbThread.Title).To(Equal(title))
				Expect(dbThread.Body).To(Equal(body))
				Expect(dbThread.UserID).To(Equal(user.ID))
			})
		})

		Context("認証トークンがない場合", func() {
			It("ステータスコード401を返す", func() {
				title := "テスト"
				body := "テストテスト"

				request := getCreateThreadRequestBodyBites(title, body)

				url := "/threads"
				w := requestAPI(http.MethodPost, url, "", request)

				Expect(w.Code).To(Equal(http.StatusUnauthorized))
			})
		})

		Context("タイトルがない場合", func() {
			It("ステータスコード400を返す", func() {
				body := "テストテスト"

				request := getCreateThreadRequestBodyBites("", body)

				url := "/threads"
				w := requestAPI(http.MethodPost, url, token, request)

				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
		})

		Context("本文がない場合", func() {
			It("ステータスコード400を返す", func() {
				title := "テスト"

				request := getCreateThreadRequestBodyBites(title, "")

				url := "/threads"
				w := requestAPI(http.MethodPost, url, token, request)

				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
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

func getThreadListResponseBody(w *httptest.ResponseRecorder) ListResponseBody {
	var body ListResponseBody
	decoder := json.NewDecoder(bytes.NewReader(w.Body.Bytes()))
	decoder.Decode(&body)

	return body
}

func getThreadDetailResponseBody(w *httptest.ResponseRecorder) DetailResponse {
	var body DetailResponse
	decoder := json.NewDecoder(bytes.NewReader(w.Body.Bytes()))
	decoder.Decode(&body)

	return body
}

func getCreateThreadRequestBodyBites(title string, body string) []byte {

	request := CreateRequest{}

	if title != "" {
		request.Title = title
	}

	if body != "" {
		request.Body = body
	}

	requestBytes, _ := json.Marshal(request)

	return requestBytes
}

func getThreadCreateResponseBody(w *httptest.ResponseRecorder) CreateResponse {
	var res CreateResponse
	decoder := json.NewDecoder(bytes.NewReader(w.Body.Bytes()))
	decoder.Decode(&res)

	return res
}
