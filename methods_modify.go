package cache

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

func (c *Cache) Clear() {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	c.data = make(cacheDataMap)
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
