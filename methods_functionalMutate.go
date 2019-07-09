package cache

import "time"

func (c *Cache) IterateClear(callback func(string, interface{}, time.Time)) {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	c.mutated = true

	for key, datum := range c.data {
		callback(key, datum.value, datum.setTime)
		delete(c.data, key)
	}
}

func (c *Cache) Filter(callback func(string, interface{}, time.Time) bool) {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	c.mutated = true

	for key, datum := range c.data {
		if retain := callback(key, datum.value, datum.setTime); !retain {
			delete(c.data, key)
		}
	}
}

func (c *Cache) Map(callback func(string, interface{}, time.Time) interface{}) {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	c.mutated = true

	for key, datum := range c.data {
		if value := callback(key, datum.value, datum.setTime); value != datum.value {
			c.data[key] = newCacheDatum(key, value)
		}
	}
}
