package main

import (
	"bbs/infra"
	"bbs/models"
)

func main() {
	infra.Init()
	db := infra.SetUpDB()

	if err := db.AutoMigrate(&models.User{}, &models.Thread{}, &models.Comment{}); err != nil {
		panic("failed to migrate database")
	}
}
