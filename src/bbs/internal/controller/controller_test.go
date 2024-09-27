package controller_test

import (
	"bbs/internal/infra"
	"bbs/internal/model"
	"bbs/internal/route"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponseBody struct {
	Token *string `json:"token"`
}

type testUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

var r *gin.Engine

var db *gorm.DB

var tmpDB *gorm.DB

var user *model.User

var token string

var contentType = "application/json"

var password = "password"

const ENV_FILE_NAME = ".env.test"

func TestBbs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Bbs Suite")
}

var _ = BeforeSuite(func() {
	infra.TestInit(getEnvTestPath())
	db = infra.SetUpDB()
	gin.SetMode(gin.TestMode)
})

func defaultBeforeEachFunc() {
	setTransaction()
	setGinRoute()
	setUserWithToken()
}

func defaultAfterEachFunc() {
	db.Rollback()
	db = tmpDB
}

func setTransaction() {
	tmpDB = db
	db = db.Begin()
}

func setGinRoute() {
	r = gin.New()
	route.SetThreadRoute(r, db)
	route.SetCommentRoute(r, db)
	route.SetAuthRoute(r, db)
}

func setUserWithToken() {
	name := "test"
	email := "exmaple@example.com"
	user = createTestUser(r, db, name, email)
	token = createTestUserToken(r, user.Email)
}

func getEnvTestPath() string {
	currentDir, _ := os.Getwd()
	return filepath.Join(currentDir, "..", "..", "configs", ENV_FILE_NAME)
}

func createTestUser(r *gin.Engine, db *gorm.DB, name string, email string) *model.User {
	u := testUser{
		Name:     name,
		Email:    email,
		Password: password,
	}

	jsonBytes, _ := json.Marshal(u)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/auth/signup", bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", contentType)
	r.ServeHTTP(w, req)

	var testUser model.User
	db.First(&testUser, "email = ?", email)

	return &testUser
}

func createTestUserToken(r *gin.Engine, email string) string {
	request := loginRequest{
		Email:    email,
		Password: password,
	}

	requestBytes, _ := json.Marshal(request)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(requestBytes))
	r.ServeHTTP(w, req)

	var body loginResponseBody
	json.Unmarshal(w.Body.Bytes(), &body)

	return *body.Token
}

func createTestThread(db *gorm.DB, userId uint, num int) []model.Thread {
	testTitle := "テストタイトル"
	testBody := "テスト本文"

	threadList := make([]model.Thread, num)

	for i := 0; i < num; i++ {
		threadList[i] = model.Thread{
			UserID: userId,
			Title:  testTitle + strconv.Itoa(i+1),
			Body:   testBody + strconv.Itoa(i+1),
		}
		db.Create(&threadList[i])
	}

	return threadList
}

func createTestComment(db *gorm.DB, userId uint, num int) []model.Comment {
	testThread := createTestThread(db, userId, 1)[0]
	testBody := "コメント本文"

	commentList := make([]model.Comment, num)

	for i := 0; i < num; i++ {
		commentList[i] = model.Comment{
			UserID:   userId,
			ThreadID: testThread.ID,
			Body:     testBody,
		}
		db.Create(&commentList[i])
	}

	return commentList
}

func requestAPI(httpMethod string, url string, authToken *string, body *[]byte) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(httpMethod, url, bytes.NewBuffer(*body))
	req.Header.Set("Content-Type", contentType)

	if authToken != nil {
		req.Header.Set("Authorization", "Bearer "+*authToken)
	}

	r.ServeHTTP(w, req)

	return w
}
