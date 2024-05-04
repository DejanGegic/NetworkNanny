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

func DefaultLimiterConf() LimiterConf {
	windowInt, err := strconv.Atoi(os.Getenv("WINDOW"))

	requestLimitInt, err := strconv.Atoi(os.Getenv("REQUEST_LIMIT"))

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
			return c.Status(429).SendString("Cool down for " + ttl.Truncate(time.Second).String() + ", until: " + time.Now().Add(ttl).Format(time.Kitchen) + ".\nRefreshes past this point will result in a ban.")
		}
		return c.Next()
	}
}

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
