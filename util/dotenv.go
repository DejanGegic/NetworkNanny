package util

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		recursiveEnvLoad()
	}
}

func recursiveEnvLoad() {

	// climb directory structure until .env file is found
	currentDir, _ := os.Getwd()

	for currentDir != "home" && currentDir != "Users" {
		if _, err := os.Stat(currentDir + "/.env"); err == nil {
			_ = godotenv.Load(currentDir + "/.env")
			return
		}
		currentDir = currentDir[:len(currentDir)-1]
	}

	panic(errors.New("no .env file found"))
}
