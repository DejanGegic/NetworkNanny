package db

import (
	"time"
)

type DbInterface interface {
	Read(key string) (string, error)
	Write(key string, value string) error
	WriteTTL(key string, value string, ttl time.Duration) error
	ReadTTL(key string) (value string, ttl time.Duration, err error)
}
