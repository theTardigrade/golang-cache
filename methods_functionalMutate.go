package cache

func (c *Cache) IterateClear(callback CallbackFunc) {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	c.mutated = true

	for key, datum := range c.data {
		callback(key, datum.value, datum.setTime)
		delete(c.data, key)
	}
}

func (c *Cache) Filter(callback CallbackFilterFunc) {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	c.mutated = true

	for key, datum := range c.data {
		if retain := callback(key, datum.value, datum.setTime); !retain {
			delete(c.data, key)
		}
	}
}

func (c *Cache) Map(callback CallbackMapFunc) {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	c.mutated = true

	for key, datum := range c.data {
		if value := callback(key, datum.value, datum.setTime); value != datum.value {
			c.data[key] = newCacheDatum(key, value)
		}
	}
}
