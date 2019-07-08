package cache

func (c *Cache) Set(key string, value interface{}) {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	c.data[key] = newCacheDatum(key, value)
	c.mutated = true
}

func (c *Cache) SetIfHasNot(key string, value interface{}) {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	if _, exists := c.data[key]; !exists {
		c.data[key] = newCacheDatum(key, value)
		c.mutated = true
	}
}

func (c *Cache) Unset(key string) {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	delete(c.data, key)
	c.mutated = true
}

func (c *Cache) Clear() {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	c.data = make(cacheDataMap)
	c.mutated = true
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
	c.mutated = true
}

func (c *Cache) Decrement(key string) {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	var count int64

	if datum, exists := c.data[key]; exists {
		countInterface := datum.value
		if countValue, ok := countInterface.(int64); ok {
			count = countValue
		}
	}

	c.data[key] = newCacheDatum(key, count-1)
	c.mutated = true
}
