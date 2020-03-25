package cache

import "time"

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

func (c *Cache) Increment(key string, updateSetTime bool) (count int64) {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	c.mutated = true

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

	return
}

func (c *Cache) Decrement(key string, updateSetTime bool) (count int64) {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	c.mutated = true

	datum, datumExists := c.data[key]
	if datumExists {
		countInterface := datum.value
		if countValue, ok := countInterface.(int64); ok {
			count = countValue - 1
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

	return
}
