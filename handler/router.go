package handler

import (
	"net/http"
	"os"
	"time"

	"example.com/middlewares/rateLimiter"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
)

type LimiterConf struct {
	RequestLimit int
	Window       time.Duration
	DbLocation   string
}

func SetUpRouter() *fiber.App {

	app := fiber.New()

	app.Use(rateLimiter.New(rateLimiter.DefaultLimiterConf()))

	app.Use(os.Getenv("BASE_URL_PATH"), filesystem.New(filesystem.Config{
		Root:   http.Dir(os.Getenv("PUBLIC_DIR")),
		Browse: getAllowBrowsing(),
	}))

	return app
}
