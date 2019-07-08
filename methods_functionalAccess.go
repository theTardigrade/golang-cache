package cache

import "time"

func (c *Cache) Iterate(callback func(string, interface{}, time.Time)) {
	defer c.mutex.RUnlock()
	c.mutex.RLock()

	for key, datum := range c.data {
		callback(key, datum.value, datum.setTime)
	}
}
