package cache

func (c *Cache) IterateClear(callback CallbackFunc) {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	for key, datum := range c.data {
		callback(key, datum.value, datum.setTime)
		c.unset(datum)
	}

	c.mutated = true
}

func (c *Cache) Filter(callback CallbackFilterFunc) {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	for key, datum := range c.data {
		if retain := callback(key, datum.value, datum.setTime); !retain {
			c.unset(datum)
		}
	}

	c.mutated = true
}

func (c *Cache) Map(callback CallbackMapFunc) {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	for key, datum := range c.data {
		if value := callback(key, datum.value, datum.setTime); value != datum.value {
			c.unset(datum)
			c.data[key] = newCacheDatum(key, value)
		}
	}

	c.mutated = true
}
