package rateLimiter

import (
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

// DefaultLimiterConf returns the default configuration for the rate limiter.
//
// This function retrieves the values of the "WINDOW", "REQUEST_LIMIT", "PERMABAN_THRESHOLD",
// and "PERMABAN_TIME" environment variables and converts them to their respective types.
// It then constructs and returns a LimiterConf struct with the retrieved values.
//
// Returns:
// - LimiterConf: the default configuration for the rate limiter.
func DefaultLimiterConf() LimiterConf {
	windowInt, err := strconv.Atoi(os.Getenv("WINDOW"))
	if err != nil {
		l.ErrorTrace(err)
		panic(err)
	}

	requestLimitInt, err := strconv.Atoi(os.Getenv("REQUEST_LIMIT"))
	if err != nil {
		l.ErrorTrace(err)
		panic(err)
	}

	permabanThresholdInt, err := strconv.Atoi(os.Getenv("PERMABAN_THRESHOLD"))
	if err != nil {
		l.ErrorTrace(err)
		panic(err)
	}
	permabanTimeInt, err := strconv.Atoi(os.Getenv("PERMABAN_TIME"))
	if err != nil {
		l.ErrorTrace(err)
		panic(err)
	}
	return LimiterConf{
		RequestLimit:      requestLimitInt,
		Window:            time.Duration(windowInt) * time.Second,
		DbLocation:        os.Getenv("DB_LOCATION"),
		LimiterName:       os.Getenv("LIMITER_NAME"),
		PermaBanTime:      time.Duration(permabanTimeInt) * time.Minute,
		PermaBanThreshold: permabanThresholdInt,
	}
}

func New(config LimiterConf) fiber.Handler {

	return func(c *fiber.Ctx) error {
		ip := c.IP()
		// timer := time.Now()
		block, ttl := checkIp(ip, config)
		// fmt.Println(time.Since(timer))

		if block == "perma" {
			return c.Status(429).SendString("PermaBanned")
		}
		if block != "" {
			return c.Status(429).SendString("Cool down for " + ttl.Truncate(time.Second).String() + ", until: " + time.Now().Add(ttl).Format(time.TimeOnly) + ".\nRefreshes past this point will result in a ban.")
		}
		return c.Next()
	}
}

// checkIp checks the IP address against the rate limiter configuration.
//
// Parameters:
// - ip: the IP address to check.
// - config: the rate limiter configuration.
// Returns:
// - block: the status of the block (perma or empty).
// - ttl: the time duration until the next request is allowed.
func checkIp(ip string, config LimiterConf) (block string, ttl time.Duration) {
	// check for PermaBan, return if true
	res, ttl, err := DB.ReadTTL("PermaBan" + ":" + ip)
	if err == nil && res != "" {
		block = "perma"
		// l.Info("Requested after PermaBan: " + ip + " by: " + config.LimiterName)
		return block, ttl
	}

	// check number of times ip has been requested in last window
	res, ttl, err = DB.ReadTTL(config.LimiterName + ":" + ip)

	// if not in db, create it
	if res == "" {
		DB.WriteTTL(config.LimiterName+":"+ip, "1", config.Window)
		return "", config.Window
	} else if err != nil {
		l.Error(err)
		return "", config.Window
	}

	resInt, _ := strconv.Atoi(res)
	if resInt >= config.RequestLimit {
		block = config.Window.String()
		if resInt >= config.PermaBanThreshold {
			DB.WriteTTL("PermaBan"+":"+ip, "1", config.PermaBanTime)
			l.Warning("PermaBanned: " + ip + " by: " + config.LimiterName + " For: " + config.PermaBanTime.String())
		}
	}
	resInt++
	DB.WriteTTL(config.LimiterName+":"+ip, strconv.Itoa(resInt), ttl)

	return block, ttl
}
