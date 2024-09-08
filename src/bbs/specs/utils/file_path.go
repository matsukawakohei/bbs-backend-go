package utils

import (
	"os"
	"path/filepath"
)

const ENV_FILE_NAME = ".env.test"

func GetEnvTestPath() string {
	currentDir, _ := os.Getwd()
	return filepath.Join(currentDir, "..", "..", "configs", ENV_FILE_NAME)
}
