package cache

import "time"

func (c *Cache) Iterate(callback func(string, interface{}, time.Time)) {
	defer c.mutex.RUnlock()
	c.mutex.RLock()

	for key, datum := range c.data {
		callback(key, datum.value, datum.setTime)
	}
}

func (c *Cache) IterateClear(callback func(string, interface{}, time.Time)) {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	for key, datum := range c.data {
		callback(key, datum.value, datum.setTime)
		delete(c.data, key)
	}
}

func (c *Cache) Filter(callback func(string, interface{}, time.Time) bool) {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	for key, datum := range c.data {
		if retain := callback(key, datum.value, datum.setTime); !retain {
			delete(c.data, key)
		}
	}
}
