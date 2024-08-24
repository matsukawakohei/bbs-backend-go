package utils

import (
	"bbs/models"

	"gorm.io/gorm"
)

func CreateTestUser(db *gorm.DB) *models.User {
	testUserName := "test"
	testUserEmail := "exmaple@example.com"
	testUserPassword := "password"
	user := models.User{
		Name:     testUserName,
		Email:    testUserEmail,
		Password: testUserPassword,
	}
	db.Create(&user)

	return &user
}
