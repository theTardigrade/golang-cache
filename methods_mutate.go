package cache

import "time"

func (c *Cache) Set(key string, value interface{}) (overwrite bool) {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	if datum, exists := c.data[key]; exists {
		c.unset(datum)
		overwrite = true
	}

	c.data[key] = newCacheDatum(key, value)
	c.hasMutatedSinceClean = true

	return
}

func (c *Cache) SetIfHasNot(key string, value interface{}) (success bool) {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	if _, exists := c.data[key]; !exists {
		c.data[key] = newCacheDatum(key, value)
		c.hasMutatedSinceClean = true
		success = true
	}

	return
}

// unset must only be called when mutex is already locked.
func (c *Cache) unset(datum *cacheDatum) {
	if c.options.UnsetPreFunc != nil {
		c.options.UnsetPreFunc(datum.key, datum.value, datum.setTime)
	}

	delete(c.data, datum.key)

	if c.options.UnsetPostFunc != nil {
		c.options.UnsetPostFunc(datum.key, datum.value, datum.setTime)
	}
}

func (c *Cache) Unset(key string) (success bool) {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	if datum, ok := c.data[key]; ok {
		c.unset(datum)
		c.hasMutatedSinceClean = true
		success = true
	}

	return
}

func (c *Cache) Clear() (overwrite bool) {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	if len(c.data) > 0 {
		for _, datum := range c.data {
			c.unset(datum)
		}

		overwrite = true
		c.hasMutatedSinceClean = true
	}

	return
}

func (c *Cache) Increment(key string, updateSetTime bool) (count int64, overwrite bool) {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	datum, datumExists := c.data[key]
	if datumExists {
		countInterface := datum.value
		if countValue, ok := countInterface.(int64); ok {
			count = countValue
		}
	}

	count++

	if datumExists {
		datum.value = count

		if updateSetTime {
			datum.setTime = time.Now()
		}
	} else {
		c.data[key] = newCacheDatum(key, count)
	}

	c.hasMutatedSinceClean = true
	overwrite = datumExists

	return
}

func (c *Cache) Decrement(key string, updateSetTime bool) (count int64, overwrite bool) {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	datum, datumExists := c.data[key]
	if datumExists {
		countInterface := datum.value
		if countValue, ok := countInterface.(int64); ok {
			count = countValue
		}
	}

	count--

	if datumExists {
		datum.value = count

		if updateSetTime {
			datum.setTime = time.Now()
		}
	} else {
		c.data[key] = newCacheDatum(key, count)
	}

	c.hasMutatedSinceClean = true
	overwrite = datumExists

	return
}
