package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponseBody struct {
	Token *string `json:"token"`
}

func CreateTestUserToken(r *gin.Engine, email string) string {
	request := loginRequest{
		Email:    email,
		Password: Password,
	}

	requestBytes, _ := json.Marshal(request)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(requestBytes))
	r.ServeHTTP(w, req)

	var body loginResponseBody
	json.Unmarshal(w.Body.Bytes(), &body)

	return *body.Token
}
