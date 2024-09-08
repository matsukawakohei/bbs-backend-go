package infra

import (
	"log"

	"github.com/joho/godotenv"
)

func Init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}
}

func TestInit(envFilePath string) {
	err := godotenv.Load(envFilePath)
	if err != nil {
		log.Fatal("error loading .env.test file", err)
	}
}
