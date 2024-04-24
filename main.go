package main

import (
	"example.com/handler"
)

func main() {

	app := handler.SetUpRouter()

	app.Listen(":8000")
}
