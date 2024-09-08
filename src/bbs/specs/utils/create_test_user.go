package utils

import (
	"bbs/internal/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type user struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func CreateTestUser(r *gin.Engine, db *gorm.DB, name string, email string) *models.User {
	u := user{
		Name:     name,
		Email:    email,
		Password: Password,
	}
	jsonBytes, _ := json.Marshal(u)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/auth/signup", bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", ContentType)
	r.ServeHTTP(w, req)

	var testUser models.User
	db.First(&testUser, "email = ?", email)

	return &testUser
}
