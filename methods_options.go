package cache

import "time"

func (c *Cache) SetMaxValues(n int) {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	c.mutated = true
	c.options.MaxValues = n
}

func (c *Cache) SetExpiryDuration(d time.Duration) {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	c.mutated = true
	c.options.ExpiryDuration = d
}
