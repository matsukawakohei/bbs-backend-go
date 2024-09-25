package controller_test

import (
	"bbs/internal/model"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type CommentCreateRequest struct {
	Body string `json:"body"`
}

type CommentCreateResponse struct {
	Comment model.Comment `json:"data"`
}

type CommentUpdateRequest struct {
	Body string `json:"body"`
}

type CommentUpdateResponse struct {
	Comment model.Comment `json:"data"`
}

var _ = Describe("CommentController", func() {
	BeforeEach(func() {
		defaultBeforeEachFunc()
	})

	AfterEach(func() {
		defaultAfterEachFunc()
	})

	Describe("コメント作成", func() {
		Context("リクエストに問題がない場合", func() {
			It("コメントを作成する", func() {
				testThreadNum := 1
				testThread := createTestThread(db, user.ID, testThreadNum)[0]

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
				testThread := createTestThread(db, user.ID, testThreadNum)[0]

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
				testThread := createTestThread(db, user.ID, testThreadNum)[0]

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

	Describe("コメント更新", func() {
		Context("リクエストが正常な場合", func() {
			It("コメントを更新する", func() {
				testCommentNum := 1
				testComment := createTestComment(db, user.ID, testCommentNum)[0]

				body := "コメント本文更新"
				request := CommentUpdateRequest{
					Body: body,
				}
				requestBytes, _ := json.Marshal(request)

				w := httptest.NewRecorder()
				url := "/threads/" + strconv.Itoa(int(testComment.ThreadID)) + "/comments/" + strconv.Itoa(int(testComment.ID))
				req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(requestBytes))
				req.Header.Set("Content-Type", contentType)
				req.Header.Set("Authorization", "Bearer "+token)
				r.ServeHTTP(w, req)

				var res CommentUpdateResponse
				decoder := json.NewDecoder(bytes.NewReader(w.Body.Bytes()))
				decoder.Decode(&res)

				Expect(err).To(BeNil())
				Expect(w.Code).To(Equal(http.StatusOK))
				Expect(res.Comment.ID).To(Equal(testComment.ID))
				Expect(res.Comment.UserID).To(Equal(user.ID))
				Expect(res.Comment.ThreadID).To(Equal(testComment.ThreadID))
				Expect(res.Comment.Body).To(Equal(body))
			})
		})

		Context("認証トークンがない場合", func() {
			It("401エラーを返す", func() {
				testCommentNum := 1
				testComment := createTestComment(db, user.ID, testCommentNum)[0]

				body := "コメント本文更新"
				request := CommentUpdateRequest{
					Body: body,
				}
				requestBytes, _ := json.Marshal(request)

				w := httptest.NewRecorder()
				url := "/threads/" + strconv.Itoa(int(testComment.ThreadID)) + "/comments/" + strconv.Itoa(int(testComment.ID))
				req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(requestBytes))
				req.Header.Set("Content-Type", contentType)
				r.ServeHTTP(w, req)

				Expect(err).To(BeNil())
				Expect(w.Code).To(Equal(http.StatusUnauthorized))
			})
		})

		Context("対象が存在しない", func() {
			It("スレッドがない場合は404エラーを返す", func() {
				testCommentNum := 1
				testComment := createTestComment(db, user.ID, testCommentNum)[0]

				body := "コメント本文更新"
				request := CommentUpdateRequest{
					Body: body,
				}
				requestBytes, _ := json.Marshal(request)

				w := httptest.NewRecorder()
				url := "/threads/" + strconv.Itoa(int(testComment.ThreadID+1)) + "/comments/" + strconv.Itoa(int(testComment.ID))
				req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(requestBytes))
				req.Header.Set("Content-Type", contentType)
				req.Header.Set("Authorization", "Bearer "+token)
				r.ServeHTTP(w, req)

				Expect(err).To(BeNil())
				Expect(w.Code).To(Equal(http.StatusNotFound))
			})

			It("コメントがない場合は404エラーを返す", func() {
				testCommentNum := 1
				testComment := createTestComment(db, user.ID, testCommentNum)[0]

				body := "コメント本文更新"
				request := CommentUpdateRequest{
					Body: body,
				}
				requestBytes, _ := json.Marshal(request)

				w := httptest.NewRecorder()
				url := "/threads/" + strconv.Itoa(int(testComment.ThreadID)) + "/comments/" + strconv.Itoa(int(testComment.ID+1))
				req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(requestBytes))
				req.Header.Set("Content-Type", contentType)
				req.Header.Set("Authorization", "Bearer "+token)
				r.ServeHTTP(w, req)

				Expect(err).To(BeNil())
				Expect(w.Code).To(Equal(http.StatusNotFound))
			})
		})

		Context("URLパラメータが文字列の場合", func() {
			It("スレッドIDが文字列の場合は400エラーを返す", func() {
				testCommentNum := 1
				testComment := createTestComment(db, user.ID, testCommentNum)[0]

				body := "コメント本文更新"
				request := CommentUpdateRequest{
					Body: body,
				}
				requestBytes, _ := json.Marshal(request)

				w := httptest.NewRecorder()
				url := "/threads/" + "aaa" + "/comments/" + strconv.Itoa(int(testComment.ID))
				req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(requestBytes))
				req.Header.Set("Content-Type", contentType)
				req.Header.Set("Authorization", "Bearer "+token)
				r.ServeHTTP(w, req)

				Expect(err).To(BeNil())
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})

			It("コメントIDが文字列の場合は400エラーを返す", func() {
				testCommentNum := 1
				testComment := createTestComment(db, user.ID, testCommentNum)[0]

				body := "コメント本文更新"
				request := CommentUpdateRequest{
					Body: body,
				}
				requestBytes, _ := json.Marshal(request)

				w := httptest.NewRecorder()
				url := "/threads/" + strconv.Itoa(int(testComment.ThreadID)) + "/comments/" + "bbb"
				req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(requestBytes))
				req.Header.Set("Content-Type", contentType)
				req.Header.Set("Authorization", "Bearer "+token)
				r.ServeHTTP(w, req)

				Expect(err).To(BeNil())
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
		})

		/** TODO 実装後にコメントアウトを外す */

		// Context("コメントの所有者と更新者が異なる", func() {
		// 	It("401エラーが返る", func() {
		// 		testCommentNum := 1
		// 		testComment := createTestComment(db, user.ID, testCommentNum)[0]

		// 		name := "test"
		// 		email := "exampleexample@example.com"
		// 		otherUser := createTestUser(r, db, name, email)
		// 		otherUserToken := createTestUserToken(r, otherUser.Email)

		// 		body := "コメント本文更新"
		// 		request := CommentUpdateRequest{
		// 			Body: body,
		// 		}
		// 		requestBytes, _ := json.Marshal(request)

		// 		w := httptest.NewRecorder()
		// 		url := "/threads/" + strconv.Itoa(int(testComment.ThreadID)) + "/comments/" + strconv.Itoa(int(testComment.ID))
		// 		req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(requestBytes))
		// 		req.Header.Set("Content-Type", contentType)
		// 		req.Header.Set("Authorization", "Bearer "+otherUserToken)
		// 		r.ServeHTTP(w, req)

		// 		var res CommentUpdateResponse
		// 		decoder := json.NewDecoder(bytes.NewReader(w.Body.Bytes()))
		// 		decoder.Decode(&res)

		// 		Expect(err).To(BeNil())
		// 		Expect(w.Code).To(Equal(http.StatusUnauthorized))
		// 	})
		// })
	})
})
