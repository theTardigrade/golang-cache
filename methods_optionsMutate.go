package cache

import "time"

func (c *Cache) SetMaxValues(n int) {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	c.setStatus(statusHasMutatedSinceCleanedFully)
	c.options.MaxValues = n

	c.startCleanWatchIfNecessary()
}

func (c *Cache) recalculateCleanInterval() {
	select {
	case c.cleanIntervalChan <- struct{}{}:
		break
	default:
		break
	}
}

func (c *Cache) SetExpiryDuration(d time.Duration) {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	c.setStatus(statusHasMutatedSinceCleanedFully)
	c.options.ExpiryDuration = d

	c.startCleanWatchIfNecessary()
	c.recalculateCleanInterval()
}

func (c *Cache) SetCleanDuration(d time.Duration) {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	c.setStatus(statusHasMutatedSinceCleanedFully)
	c.options.CleanDuration = d

	c.recalculateCleanInterval()
}

func (c *Cache) SetCleanMaxValuesPerSweep(n int) {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	c.setStatus(statusHasMutatedSinceCleanedFully)
	c.options.CleanMaxValuesPerSweep = n
}

func (c *Cache) SetUnsetPreFunc(p CallbackFunc) {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	c.setStatus(statusHasMutatedSinceCleanedFully)
	c.options.UnsetPreFunc = p
}

func (c *Cache) SetUnsetPostFunc(p CallbackFunc) {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	c.setStatus(statusHasMutatedSinceCleanedFully)
	c.options.UnsetPostFunc = p
}
