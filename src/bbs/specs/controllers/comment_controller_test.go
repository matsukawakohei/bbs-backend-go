package controllers_test

import (
	"bbs/models"
	"bbs/specs/utils"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type CommentCreateRequest struct {
	Body string `json:"body"`
}

type CommentCreateResponse struct {
	Comment models.Comment `json:"data"`
}

func TestComment(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Comment Suite")
}

var r = utils.R

var db = utils.Db

var user = utils.User

var token = utils.Token

var contentType = utils.ContentType

var _ = Describe("CommentController", func() {
	AfterEach(func() {
		db.Where("id > ?", 0).Unscoped().Delete(&models.Thread{})
		db.Where("id > ?", 0).Unscoped().Delete(&models.Comment{})
	})

	Describe("コメント作成", func() {
		Context("リクエストに問題がない場合", func() {
			It("コメントを作成する", func() {
				testThreadNum := 1
				testThread := utils.CreateTestThread(db, user.ID, testThreadNum)[0]

				body := "コメント本文"
				request := CommentCreateRequest{
					Body: body,
				}
				requestBytes, _ := json.Marshal(request)

				w := httptest.NewRecorder()
				url := "/threads/" + strconv.Itoa(int(testThread.ID)) + "/comments"
				req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBytes))
				req.Header.Set("Content-Type", contentType)
				req.Header.Set("Authorization", "Bearer "+token)
				r.ServeHTTP(w, req)

				var res CommentCreateResponse
				decoder := json.NewDecoder(bytes.NewReader(w.Body.Bytes()))
				decoder.Decode(&res)

				Expect(err).To(BeNil())
				Expect(w.Code).To(Equal(http.StatusCreated))
				Expect(res.Comment.ID).NotTo(BeNil())
				Expect(res.Comment.UserID).To(Equal(user.ID))
				Expect(res.Comment.ThreadID).To(Equal(testThread.ID))
				Expect(res.Comment.Body).To(Equal(body))
			})
		})

		Context("ログインしていない場合", func() {
			It("401エラーが返る", func() {
				testThreadNum := 1
				testThread := utils.CreateTestThread(db, user.ID, testThreadNum)[0]

				body := "コメント本文"
				request := CommentCreateRequest{
					Body: body,
				}
				requestBytes, _ := json.Marshal(request)

				w := httptest.NewRecorder()
				url := "/threads/" + strconv.Itoa(int(testThread.ID)) + "/comments"
				req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBytes))
				req.Header.Set("Content-Type", contentType)
				r.ServeHTTP(w, req)

				Expect(err).To(BeNil())
				Expect(w.Code).To(Equal(http.StatusUnauthorized))
			})
		})

		Context("threadIDが数値ではない場合", func() {
			It("400エラーが返る", func() {
				body := "コメント本文"
				request := CommentCreateRequest{
					Body: body,
				}
				requestBytes, _ := json.Marshal(request)

				w := httptest.NewRecorder()
				url := "/threads/aaa/comments"
				req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBytes))
				req.Header.Set("Content-Type", contentType)
				req.Header.Set("Authorization", "Bearer "+token)
				r.ServeHTTP(w, req)

				Expect(err).To(BeNil())
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
		})

		Context("リクエストパラメーターにBodyがない場合", func() {
			It("400エラーが返る", func() {
				testThreadNum := 1
				testThread := utils.CreateTestThread(db, user.ID, testThreadNum)[0]

				w := httptest.NewRecorder()
				url := "/threads/" + strconv.Itoa(int(testThread.ID)) + "/comments"
				req, err := http.NewRequest(http.MethodPost, url, nil)
				req.Header.Set("Content-Type", contentType)
				req.Header.Set("Authorization", "Bearer "+token)
				r.ServeHTTP(w, req)

				Expect(err).To(BeNil())
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
		})

		Context("スレッドが存在しない場合", func() {
			It("404エラーが返る", func() {
				body := "コメント本文"
				request := CommentCreateRequest{
					Body: body,
				}
				requestBytes, _ := json.Marshal(request)

				w := httptest.NewRecorder()
				url := "/threads/1/comments"
				req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBytes))
				req.Header.Set("Content-Type", contentType)
				req.Header.Set("Authorization", "Bearer "+token)
				r.ServeHTTP(w, req)

				Expect(err).To(BeNil())
				Expect(w.Code).To(Equal(http.StatusNotFound))
			})
		})
	})
})