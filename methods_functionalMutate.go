package cache

func (c *Cache) IterateClear(callback CallbackFunc) (overwrite bool) {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	if len(c.data) > 0 {
		for key, datum := range c.data {
			callback(key, datum.value, datum.setTime)
			c.unset(datum)
		}

		c.mutated = true
		overwrite = true
	}

	return
}

func (c *Cache) Filter(callback CallbackFilterFunc) (overwrite bool) {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	for key, datum := range c.data {
		if retain := callback(key, datum.value, datum.setTime); !retain {
			c.unset(datum)
			overwrite = true
		}
	}

	if !c.mutated {
		c.mutated = overwrite
	}

	return
}

func (c *Cache) Map(callback CallbackMapFunc) (overwrite bool) {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	for key, datum := range c.data {
		if value := callback(key, datum.value, datum.setTime); value != datum.value {
			c.unset(datum)
			c.data[key] = newCacheDatum(key, value)
			overwrite = true
		}
	}

	if !c.mutated {
		c.mutated = overwrite
	}

	return
}
