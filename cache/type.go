package cache

import (
	"time"
)

type AdapterCache interface {
	Connect() error
	DisConnect() error

	Get(key string) (string, error)
	Set(key string, value string, expire time.Duration) error
	Del(key string) error
	Expire(key string, dur time.Duration) error
	Exists(key string) bool

	HGetAll(key string) (map[string]string, error)
	HGet(key, field string) (string, error)
	HSet(key, field, value string) error
	HDel(key string, field string) error
	HExists(key, field string) (bool, error)
}
