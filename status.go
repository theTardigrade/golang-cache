package cache

type status uint8

const (
	statusHasMutatedSinceCleanedFully = 1 << iota
	statusHasCleanWatchStarted
)

func (c *Cache) setStatus(s status) {
	c.status |= s
}

func (c *Cache) unsetStatus(s status) {
	c.status &= ^s
}

func (c *Cache) hasStatus(s status) bool {
	return (c.status & s) != 0
}
