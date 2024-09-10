package timer

import (
	"time"

	"example.com/l"
	"github.com/gofiber/fiber/v2"
)

// Fiber middleware function that prints the time it took to process the request
func Timer() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		duration := time.Since(start)
		l.Z.Info().Str("IP", c.IP()).Str("method", c.Method()).Str("path", c.Path()).Dur("duration", duration).Send()
		return err
	}
}
