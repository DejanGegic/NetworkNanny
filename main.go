package main

import (
	"os"

	"example.com/handler"
)

func main() {

	app := handler.SetUpRouter()
	app.Listen(os.Getenv("PORT"))
}
