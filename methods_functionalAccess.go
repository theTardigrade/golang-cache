package cache

func (c *Cache) Iterate(callback CallbackFunc) {
	defer c.mutex.RUnlock()
	c.mutex.RLock()

	for key, datum := range c.data {
		callback(key, datum.value, datum.setTime)
	}
}
