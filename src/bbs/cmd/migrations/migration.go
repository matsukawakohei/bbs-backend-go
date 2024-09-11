package main

import (
	"bbs/internal/infra"
	"bbs/internal/model"
)

func main() {
	infra.Init()
	db := infra.SetUpDB()

	if err := db.AutoMigrate(&model.User{}, &model.Thread{}, &model.Comment{}); err != nil {
		panic("failed to migrate database")
	}
}
