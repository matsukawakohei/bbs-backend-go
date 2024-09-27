package controller_test

import (
	"bbs/internal/model"
	"bytes"
	"encoding/json"
	"net/http"
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
		Context("リクエストが正常な場合", func() {
			It("ステータスコード201が返る", func() {
				testThreadNum := 1
				testThread := createTestThread(db, user.ID, testThreadNum)[0]

				body := "コメント本文"
				requestBytes := getCreateCommentRequestBodyBites(body)

				url := "/threads/" + strconv.Itoa(int(testThread.ID)) + "/comments"
				w := requestAPI(http.MethodPost, url, &token, &requestBytes)

				Expect(w.Code).To(Equal(http.StatusCreated))
			})

			It("作成したコメントが返る", func() {
				testThreadNum := 1
				testThread := createTestThread(db, user.ID, testThreadNum)[0]

				body := "コメント本文"
				requestBytes := getCreateCommentRequestBodyBites(body)

				url := "/threads/" + strconv.Itoa(int(testThread.ID)) + "/comments"
				w := requestAPI(http.MethodPost, url, &token, &requestBytes)

				var res CommentCreateResponse
				decoder := json.NewDecoder(bytes.NewReader(w.Body.Bytes()))
				decoder.Decode(&res)

				Expect(res.Comment.ID).NotTo(BeNil())
				Expect(res.Comment.UserID).To(Equal(user.ID))
				Expect(res.Comment.ThreadID).To(Equal(testThread.ID))
				Expect(res.Comment.Body).To(Equal(body))
			})

			It("DBに作成したコメントが登録されている", func() {
				testThreadNum := 1
				testThread := createTestThread(db, user.ID, testThreadNum)[0]

				body := "コメント本文"
				requestBytes := getCreateCommentRequestBodyBites(body)

				url := "/threads/" + strconv.Itoa(int(testThread.ID)) + "/comments"
				requestAPI(http.MethodPost, url, &token, &requestBytes)

				var dbComment model.Comment
				result := db.First(&dbComment)

				Expect(result.Error).To(BeNil())
				Expect(dbComment.ID).NotTo(BeNil())
				Expect(dbComment.UserID).To(Equal(user.ID))
				Expect(dbComment.ThreadID).To(Equal(testThread.ID))
				Expect(dbComment.Body).To(Equal(body))
			})
		})

		Context("ログインしていない場合", func() {
			It("401エラーが返る", func() {
				testThreadNum := 1
				testThread := createTestThread(db, user.ID, testThreadNum)[0]

				body := "コメント本文"
				requestBytes := getCreateCommentRequestBodyBites(body)

				url := "/threads/" + strconv.Itoa(int(testThread.ID)) + "/comments"
				w := requestAPI(http.MethodPost, url, nil, &requestBytes)

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

				url := "/threads/aaa/comments"
				w := requestAPI(http.MethodPost, url, &token, &requestBytes)

				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
		})

		Context("リクエストパラメーターにBodyがない場合", func() {
			It("400エラーが返る", func() {
				testThreadNum := 1
				testThread := createTestThread(db, user.ID, testThreadNum)[0]

				url := "/threads/" + strconv.Itoa(int(testThread.ID)) + "/comments"
				w := requestAPI(http.MethodPost, url, &token, nil)

				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
		})

		Context("スレッドが存在しない場合", func() {
			It("404エラーが返る", func() {
				body := "コメント本文"
				requestBytes := getCreateCommentRequestBodyBites(body)

				url := "/threads/1/comments"
				w := requestAPI(http.MethodPost, url, &token, &requestBytes)

				Expect(w.Code).To(Equal(http.StatusNotFound))
			})
		})
	})

	Describe("コメント更新", func() {
		Context("リクエストが正常な場合", func() {
			It("ステータスコード200が返る", func() {
				testCommentNum := 1
				testComment := createTestComment(db, user.ID, testCommentNum)[0]

				body := "コメント本文更新"
				requestBytes := getUpdateCommentRequestBodyBites(body)

				url := "/threads/" + strconv.Itoa(int(testComment.ThreadID)) + "/comments/" + strconv.Itoa(int(testComment.ID))
				w := requestAPI(http.MethodPut, url, &token, &requestBytes)

				Expect(w.Code).To(Equal(http.StatusOK))
			})

			It("更新後のコメントを返す", func() {
				testCommentNum := 1
				testComment := createTestComment(db, user.ID, testCommentNum)[0]

				body := "コメント本文更新"
				requestBytes := getUpdateCommentRequestBodyBites(body)

				url := "/threads/" + strconv.Itoa(int(testComment.ThreadID)) + "/comments/" + strconv.Itoa(int(testComment.ID))
				w := requestAPI(http.MethodPut, url, &token, &requestBytes)

				var res CommentUpdateResponse
				decoder := json.NewDecoder(bytes.NewReader(w.Body.Bytes()))
				decoder.Decode(&res)

				Expect(res.Comment.ID).To(Equal(testComment.ID))
				Expect(res.Comment.UserID).To(Equal(user.ID))
				Expect(res.Comment.ThreadID).To(Equal(testComment.ThreadID))
				Expect(res.Comment.Body).To(Equal(body))
			})

			It("DBのコメントが更新されている", func() {
				testCommentNum := 1
				testComment := createTestComment(db, user.ID, testCommentNum)[0]

				body := "コメント本文更新"
				requestBytes := getUpdateCommentRequestBodyBites(body)

				url := "/threads/" + strconv.Itoa(int(testComment.ThreadID)) + "/comments/" + strconv.Itoa(int(testComment.ID))
				requestAPI(http.MethodPut, url, &token, &requestBytes)

				var dbComment model.Comment
				result := db.First(&dbComment, "id = ?", testComment.ID)

				Expect(result.Error).To(BeNil())
				Expect(dbComment.ID).To(Equal(testComment.ID))
				Expect(dbComment.UserID).To(Equal(user.ID))
				Expect(dbComment.ThreadID).To(Equal(testComment.ThreadID))
				Expect(dbComment.Body).To(Equal(body))
			})
		})

		Context("認証トークンがない場合", func() {
			It("401エラーを返す", func() {
				testCommentNum := 1
				testComment := createTestComment(db, user.ID, testCommentNum)[0]

				body := "コメント本文更新"
				requestBytes := getUpdateCommentRequestBodyBites(body)

				url := "/threads/" + strconv.Itoa(int(testComment.ThreadID)) + "/comments/" + strconv.Itoa(int(testComment.ID))
				w := requestAPI(http.MethodPut, url, nil, &requestBytes)

				Expect(w.Code).To(Equal(http.StatusUnauthorized))
			})
		})

		Context("対象が存在しない", func() {
			It("スレッドがない場合は404エラーを返す", func() {
				testCommentNum := 1
				testComment := createTestComment(db, user.ID, testCommentNum)[0]

				body := "コメント本文更新"
				requestBytes := getUpdateCommentRequestBodyBites(body)

				url := "/threads/" + strconv.Itoa(int(testComment.ThreadID+1)) + "/comments/" + strconv.Itoa(int(testComment.ID))
				w := requestAPI(http.MethodPut, url, &token, &requestBytes)

				Expect(w.Code).To(Equal(http.StatusNotFound))
			})

			It("コメントがない場合は404エラーを返す", func() {
				testCommentNum := 1
				testComment := createTestComment(db, user.ID, testCommentNum)[0]

				body := "コメント本文更新"
				requestBytes := getUpdateCommentRequestBodyBites(body)

				url := "/threads/" + strconv.Itoa(int(testComment.ThreadID)) + "/comments/" + strconv.Itoa(int(testComment.ID+1))
				w := requestAPI(http.MethodPut, url, &token, &requestBytes)

				Expect(w.Code).To(Equal(http.StatusNotFound))
			})
		})

		Context("URLパラメータが文字列の場合", func() {
			It("スレッドIDが文字列の場合は400エラーを返す", func() {
				testCommentNum := 1
				testComment := createTestComment(db, user.ID, testCommentNum)[0]

				body := "コメント本文更新"
				requestBytes := getUpdateCommentRequestBodyBites(body)

				url := "/threads/" + "aaa" + "/comments/" + strconv.Itoa(int(testComment.ID))
				w := requestAPI(http.MethodPut, url, &token, &requestBytes)

				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})

			It("コメントIDが文字列の場合は400エラーを返す", func() {
				testCommentNum := 1
				testComment := createTestComment(db, user.ID, testCommentNum)[0]

				body := "コメント本文更新"
				requestBytes := getUpdateCommentRequestBodyBites(body)

				url := "/threads/" + strconv.Itoa(int(testComment.ThreadID)) + "/comments/" + "bbb"
				w := requestAPI(http.MethodPut, url, &token, &requestBytes)

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

func getCreateCommentRequestBodyBites(body string) []byte {
	request := CommentCreateRequest{
		Body: body,
	}
	requestBytes, _ := json.Marshal(request)

	return requestBytes
}

func getUpdateCommentRequestBodyBites(body string) []byte {
	request := CommentUpdateRequest{
		Body: body,
	}
	requestBytes, _ := json.Marshal(request)

	return requestBytes
}
