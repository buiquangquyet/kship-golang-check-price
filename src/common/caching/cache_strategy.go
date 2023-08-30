package caching

import (
	"time"
)

type CacheStrategy interface {
	Get(key string) (interface{}, bool)
	Set(key string, data interface{}, ttl time.Duration)
}
