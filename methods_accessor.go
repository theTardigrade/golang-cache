package cache

import "strings"

func (c *Cache) Set(key string, value interface{}) {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	c.data[key] = newCacheDatum(key, value)
}

func (c *Cache) Unset(key string) {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	delete(c.data, key)
}

func (c *Cache) Get(key string) (interface{}, bool) {
	defer c.mutex.RUnlock()
	c.mutex.RLock()

	datum, ok := c.data[key]
	if !ok {
		return nil, false
	}

	return datum.value, true
}

func (c *Cache) MustGet(key string) interface{} {
	value, ok := c.Get(key)
	if !ok {
		panic(ErrNotFound)
	}

	return value
}

func (c *Cache) Has(key string) bool {
	defer c.mutex.RUnlock()
	c.mutex.RLock()

	_, ok := c.data[key]

	return ok
}

func (c *Cache) Increment(key string) {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	var count int64

	if datum, exists := c.data[key]; exists {
		countInterface := datum.value
		if countValue, ok := countInterface.(int64); ok {
			count = countValue
		}
	}

	c.data[key] = newCacheDatum(key, count+1)
}

func (c *Cache) Clear() {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	c.data = make(cacheDataMap)
}

func (c *Cache) Len() int {
	defer c.mutex.RUnlock()
	c.mutex.RLock()

	return len(c.data)
}

func (c *Cache) String() string {
	defer c.mutex.RUnlock()
	c.mutex.RLock()

	var builder strings.Builder

	builder.WriteByte('[')

	var i int
	for key := range c.data {
		if i++; i > 1 {
			builder.WriteString(", ")
		}

		builder.WriteString(key)
	}

	builder.WriteByte(']')

	return builder.String()
}
