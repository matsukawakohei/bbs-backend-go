package utils

import (
	"bbs/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var Password = "password"

type user struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func CreateTestUser(r *gin.Engine, db *gorm.DB) *models.User {
	name := "test"
	email := "exmaple@example.com"

	u := user{
		Name:     name,
		Email:    email,
		Password: Password,
	}
	jsonBytes, _ := json.Marshal(u)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/auth/signup", bytes.NewBuffer(jsonBytes))
	r.ServeHTTP(w, req)

	var testUser models.User
	db.First(&testUser, "email = ?", email)

	return &testUser
}
