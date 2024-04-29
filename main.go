package main

import (
	"os"

	"example.com/db/redis"
	"example.com/handler"
	util "example.com/util"
)

func main() {

	util.LoadEnv()
	redis.Init()
	app := handler.SetUpRouter()
	app.Listen(os.Getenv("PORT"))
}
