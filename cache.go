package cache

import (
	"sync"
	"time"
)

type cacheDatum struct {
	key     string
	value   interface{}
	setTime time.Time
}

type cacheDataMap map[string]*cacheDatum

type Cache struct {
	data    cacheDataMap
	mutex   sync.RWMutex
	mutated bool
	options Options
}

type CallbackFunc func(string, interface{}, time.Time)
type CallbackFilterFunc func(string, interface{}, time.Time) bool
type CallbackMapFunc func(string, interface{}, time.Time) interface{}

func NewInfiniteCache() (cache *Cache) {
	cache = &Cache{}

	go cache.watch()

	return cache
}

func NewCache(expiryDuration time.Duration, maxValues int) *Cache {
	options := Options{
		ExpiryDuration: expiryDuration,
		MaxValues:      maxValues,
	}

	return NewCacheWithOptions(options)
}

func NewCacheWithOptions(options Options) *Cache {
	cache := NewInfiniteCache()

	cache.options = options

	return cache
}

func newCacheDatum(key string, value interface{}) *cacheDatum {
	return &cacheDatum{
		key:     key,
		value:   value,
		setTime: time.Now(),
	}
}
