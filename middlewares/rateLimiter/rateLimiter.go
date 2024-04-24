package rateLimiter

import (
	"fmt"
	"strconv"
	"time"

	"example.com/db"
	"example.com/l"
	"github.com/gofiber/fiber/v2"
)

type LimiterConf struct {
	RequestLimit int
	Window       time.Duration
	DbLocation   string
	LimiterName  string
}

func NewConfig() LimiterConf {
	return LimiterConf{
		RequestLimit: 5,
		Window:       60 * time.Second,
		DbLocation:   "rateLimiter.db",
		LimiterName:  "60s",
	}
}

func New(config LimiterConf) fiber.Handler {

	return func(c *fiber.Ctx) error {
		timer := time.Now()
		ip := c.IP()
		block := checkIp(ip, config)

		fmt.Println(time.Since(timer))

		if block == "perma" {
			return c.Status(429).SendString("PermaBanned")
		}
		if block != "" {
			return c.Status(429).SendString("Cool down for " + block + ". Refreshing will reset the cooldown.")
		}
		return c.Next()
	}
}

func checkIp(ip string, config LimiterConf) (block string) {
	// check for PermaBan
	permaRes, err := db.Read("PermaBan" + "|" + ip)
	if err == nil && permaRes != "" {
		block = "perma"
		return block
	}

	// TODO check if it's in the database, if not, add
	res, err := db.Read(config.LimiterName + "|" + ip)
	if err != nil {
		fmt.Println(err)
		db.WriteTTL(config.LimiterName+"|"+ip, "1", config.Window)
	}

	resInt, _ := strconv.Atoi(res)

	if resInt >= config.RequestLimit {
		block = config.Window.String()
		if resInt >= 30 {
			db.WriteTTL("PermaBan"+"|"+ip, "1", time.Hour*24)
			l.Warning("PermaBanned: " + ip)
		}
	}
	resInt++
	db.WriteTTL(config.LimiterName+"|"+ip, strconv.Itoa(resInt), config.Window)

	return block
}
