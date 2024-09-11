package utils

import (
	"bbs/internal/model"
	"strconv"

	"gorm.io/gorm"
)

func CreateTestThread(db *gorm.DB, userId uint, num int) []model.Thread {
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
