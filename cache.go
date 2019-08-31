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
	data           cacheDataMap
	mutex          sync.RWMutex
	mutated        bool
	expiryDuration time.Duration
	maxValues      int
}

func NewCache(expiryDuration time.Duration, maxValues int) *Cache {
	cache := Cache{
		data:           make(cacheDataMap),
		expiryDuration: expiryDuration, // -1 == infinite
		maxValues:      maxValues,      // -1 == infinite
	}

	if expiryDuration >= 0 || maxValues >= 0 {
		go cache.watch()
	}

	return &cache
}

func newCacheDatum(key string, value interface{}) *cacheDatum {
	return &cacheDatum{
		key:     key,
		value:   value,
		setTime: time.Now(),
	}
}
