package utils

import (
	"bbs/models"
	"strconv"

	"gorm.io/gorm"
)

func CreateTestThread(db *gorm.DB, userId uint, num int) {
	testTitle := "テストタイトル"
	testBody := "テスト本文"

	for i := 0; i < num; i++ {
		thread := models.Thread{
			UserID: userId,
			Title:  testTitle + strconv.Itoa(i+1),
			Body:   testBody + strconv.Itoa(i+1),
		}
		db.Create(&thread)
	}
}
