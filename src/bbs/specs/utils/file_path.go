package utils

import (
	"os"
	"path/filepath"
)

const ENV_FILE_NAME = ".env.test"

func GetEnvTestPath() string {
	return getFilePath(ENV_FILE_NAME)
}

func getFilePath(fileName string) string {
	currentDir, err := os.Getwd()
	if err != nil {
		return ""
	}

	for {
		_, err := os.ReadFile(filepath.Join(currentDir, fileName))
		if os.IsNotExist(err) {
			if currentDir == filepath.Dir(currentDir) {
				return ""
			}
			currentDir = filepath.Dir(currentDir)

			continue
		} else if err != nil {
			return ""
		}
		break
	}
	return filepath.Join(currentDir, fileName)
}
