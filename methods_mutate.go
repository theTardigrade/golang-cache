package cache

import (
	"reflect"
	"time"
)

func (c *Cache) Set(key string, value interface{}) (overwrite bool) {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	if datum, exists := c.data[key]; exists {
		var valuesMatch bool

		oldValue := datum.value
		valueType := reflect.TypeOf(value)
		oldValueType := reflect.TypeOf(oldValue)

		if valueType.Kind() == oldValueType.Kind() && valueType.Comparable() {
			overwrite = true
			valuesMatch = (value == oldValue)
		}

		if !valuesMatch {
			c.unset(datum)
		}

		overwrite = true
	}

	c.data[key] = newCacheDatum(key, value)
	c.setStatus(statusHasMutatedSinceCleanedFully)

	return
}

func (c *Cache) SetIfHasNot(key string, value interface{}) (success bool) {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	if _, exists := c.data[key]; !exists {
		c.data[key] = newCacheDatum(key, value)
		c.setStatus(statusHasMutatedSinceCleanedFully)
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
		c.setStatus(statusHasMutatedSinceCleanedFully)
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

		c.setStatus(statusHasMutatedSinceCleanedFully)
		overwrite = true
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

	c.setStatus(statusHasMutatedSinceCleanedFully)
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

	c.setStatus(statusHasMutatedSinceCleanedFully)
	overwrite = datumExists

	return
}
