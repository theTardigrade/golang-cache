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
	data              cacheDataMap
	mutex             sync.RWMutex
	options           Options
	cleanIntervalChan chan struct{}
	status            status
}

type CallbackFunc func(key string, value interface{}, setTime time.Time)
type CallbackFilterFunc func(key string, value interface{}, setTime time.Time) (retain bool)
type CallbackMapFunc func(key string, value interface{}, setTime time.Time) (mappedValue interface{})

func NewInfiniteCache() (cache *Cache) {
	cache = &Cache{
		data:              make(cacheDataMap),
		cleanIntervalChan: make(chan struct{}),
	}

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

	cache.startWatchIfNecessary()

	return cache
}

func newCacheDatum(key string, value interface{}) *cacheDatum {
	return &cacheDatum{
		key:     key,
		value:   value,
		setTime: time.Now(),
	}
}
