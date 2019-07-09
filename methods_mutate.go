package cache

func (c *Cache) Set(key string, value interface{}) {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	c.mutated = true
	c.data[key] = newCacheDatum(key, value)
}

func (c *Cache) SetIfHasNot(key string, value interface{}) {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	if _, exists := c.data[key]; !exists {
		c.mutated = true
		c.data[key] = newCacheDatum(key, value)
	}
}

func (c *Cache) Unset(key string) {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	c.mutated = true
	delete(c.data, key)
}

func (c *Cache) Clear() {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	c.mutated = true
	c.data = make(cacheDataMap)
}

func (c *Cache) Increment(key string) {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	c.mutated = true

	var count int64

	if datum, exists := c.data[key]; exists {
		countInterface := datum.value
		if countValue, ok := countInterface.(int64); ok {
			count = countValue
		}
	}

	c.data[key] = newCacheDatum(key, count+1)
}

func (c *Cache) Decrement(key string) {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	c.mutated = true

	var count int64

	if datum, exists := c.data[key]; exists {
		countInterface := datum.value
		if countValue, ok := countInterface.(int64); ok {
			count = countValue
		}
	}

	c.data[key] = newCacheDatum(key, count-1)
}
