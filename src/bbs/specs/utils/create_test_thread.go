package utils

import (
	"bbs/models"
	"strconv"

	"gorm.io/gorm"
)

func CreateTestThread(db *gorm.DB, userId uint, num int) []models.Thread {
	testTitle := "テストタイトル"
	testBody := "テスト本文"

	threadList := make([]models.Thread, num)

	for i := 0; i < num; i++ {
		threadList[i] = models.Thread{
			UserID: userId,
			Title:  testTitle + strconv.Itoa(i+1),
			Body:   testBody + strconv.Itoa(i+1),
		}
		db.Create(&threadList[i])
	}

	return threadList
}
