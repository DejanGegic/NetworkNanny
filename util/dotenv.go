package util

import (
	"errors"
	"os"
	"strconv"

	"example.com/l"
	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		recursiveEnvLoad()
	}
	envDefaults()
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

func envDefaults() {
	// WINDOW
	windowInt, err := strconv.Atoi(os.Getenv("WINDOW"))
	if (err != nil) || (windowInt < 1) {
		l.Error(errors.New("environment variable 'WINDOW' must be an integer greater than 0. Please check your .env file"))
		panic(err)
	}
	// REQUEST_LIMIT
	requestLimitInt, err := strconv.Atoi(os.Getenv("REQUEST_LIMIT"))
	if (err != nil) || (requestLimitInt < 1) {
		l.Error(errors.New("environment variable 'REQUEST_LIMIT' must be an integer greater than 0. Please check your .env file"))
		panic(err)
	}
	// LIMITER_NAME
	if os.Getenv("LIMITER_NAME") == "" {
		os.Setenv("LIMITER_NAME", "60s")
	}
	// DB_LOCATION
	if os.Getenv("DB_LOCATION") == "" {
		os.Setenv("DB_LOCATION", "badger")
	}
	// PERMABAN_THRESHOLD
	permabanThresholdInt, err := strconv.Atoi(os.Getenv("PERMABAN_THRESHOLD"))
	if (err != nil) || (permabanThresholdInt < 1) {
		os.Setenv("PERMABAN_THRESHOLD", "10")
	}
	// PERMABAN_TIME
	permabanTimeInt, err := strconv.Atoi(os.Getenv("PERMABAN_TIME"))
	if (err != nil) || (permabanTimeInt < 1) {
		os.Setenv("PERMABAN_TIME", "1440")
	}
	// DB_TYPE
	if os.Getenv("DB_TYPE") == "" {
		os.Setenv("DB_TYPE", "badger")
	}
	// BASE_URL_PATH
	if os.Getenv("BASE_URL_PATH") == "" {
		os.Setenv("BASE_URL_PATH", "/")
	}
	// PUBLIC_DIR
	if os.Getenv("PUBLIC_DIR") == "" {
		os.Setenv("PUBLIC_DIR", "public")
	}
	// PORT
	if os.Getenv("PORT") == "" {
		os.Setenv("PORT", "8000")
	}

	l.Info("Environment variables loaded." + " | " +
		"WINDOW: " + os.Getenv("WINDOW") + " | " +
		"REQUEST_LIMIT: " + os.Getenv("REQUEST_LIMIT") + " | " +
		"LIMITER_NAME: " + os.Getenv("LIMITER_NAME") + " | " +
		"DB_LOCATION: " + os.Getenv("DB_LOCATION") + " | " +
		"PERMABAN_THRESHOLD: " + os.Getenv("PERMABAN_THRESHOLD") + " | " +
		"PERMABAN_TIME: " + os.Getenv("PERMABAN_TIME") + " | " +
		"DB_TYPE: " + os.Getenv("DB_TYPE") + " | " +
		"BASE_URL_PATH: " + os.Getenv("BASE_URL_PATH") + " | " +
		"PUBLIC_DIR: " + os.Getenv("PUBLIC_DIR") + " | " +
		"PORT: " + os.Getenv("PORT"))
}
