package rateLimiter

import (
	"errors"
	"os"
	"strconv"
	"time"

	"example.com/db"
	"example.com/l"
	"github.com/gofiber/fiber/v2"
)

type LimiterConf struct {
	RequestLimit      int
	Window            time.Duration
	DbLocation        string // Only applicable for badger
	LimiterName       string
	PermaBanTime      time.Duration
	PermaBanThreshold int
}

var DB db.DbInterface = db.InitDB()

func DefaultLimiterConf() LimiterConf {
	windowInt, err := strconv.Atoi(os.Getenv("WINDOW"))
	if (err != nil) || (windowInt < 1) {
		l.Error(errors.New("Environment variable 'WINDOW' must be an integer greater than 0. Please check your .env file."))
		panic(err)
	}
	requestLimitInt, err := strconv.Atoi(os.Getenv("REQUEST_LIMIT"))
	if (err != nil) || (requestLimitInt < 1) {
		l.Error(errors.New("Environment variable 'REQUEST_LIMIT' must be an integer greater than 0. Please check your .env file."))
		panic(err)
	}
	name := os.Getenv("LIMITER_NAME")
	if name == "" {
		name = "60s"
	}
	dbLocation := os.Getenv("DB_LOCATION")
	if dbLocation == "" {
		dbLocation = "badger"
	}
	permabanThresholdInt, err := strconv.Atoi(os.Getenv("PERMABAN_THRESHOLD"))
	if (err != nil) || (permabanThresholdInt < 1) {
		permabanThresholdInt = 10
	}
	permabanTimeInt, err := strconv.Atoi(os.Getenv("PERMABAN_TIME"))
	if (err != nil) || (permabanTimeInt < 1) {
		permabanTimeInt = 1440
	}
	return LimiterConf{
		RequestLimit:      requestLimitInt,
		Window:            time.Duration(windowInt) * time.Second,
		DbLocation:        dbLocation,
		LimiterName:       name,
		PermaBanTime:      time.Duration(permabanTimeInt) * time.Minute,
		PermaBanThreshold: permabanThresholdInt,
	}
}

func New(config LimiterConf) fiber.Handler {

	return func(c *fiber.Ctx) error {
		ip := c.IP()
		block := checkIp(ip, config)

		if block == "perma" {
			return c.Status(429).SendString("PermaBanned")
		}
		if block != "" {
			return c.Status(429).SendString("Cool down for " + block + ". Refreshing will reset the cooldown. \nRefreshes past this point will result in a ban.")
		}
		return c.Next()
	}
}

func checkIp(ip string, config LimiterConf) (block string) {
	// check for PermaBan, return if true
	res, err := DB.Read("PermaBan" + "|" + ip)
	if err == nil && res != "" {
		block = "perma"
		return block
	}

	// check number of times ip has been requested in last window
	res, err = DB.Read(config.LimiterName + "|" + ip)
	if err != nil {
		// "Key not found" is expected.
		// Any other error is logged accordingly
		if err.Error() == "Key not found" {
			DB.WriteTTL(config.LimiterName+"|"+ip, "1", config.Window)
		} else {
			l.Error(err)
		}
	}

	resInt, _ := strconv.Atoi(res)
	if resInt >= config.RequestLimit {
		block = config.Window.String()
		if resInt >= config.PermaBanThreshold {
			DB.WriteTTL("PermaBan"+"|"+ip, "1", config.PermaBanTime)
			l.Warning("PermaBanned: " + ip + " by: " + config.LimiterName + " For: " + config.PermaBanTime.String())
		}
	}
	resInt++
	DB.WriteTTL(config.LimiterName+"|"+ip, strconv.Itoa(resInt), config.Window)

	return block
}
