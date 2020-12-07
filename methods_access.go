package cache

import "strings"

func (c *Cache) Get(key string) (value interface{}, found bool) {
	defer c.mutex.RUnlock()
	c.mutex.RLock()

	datum, found := c.data[key]
	if found {
		value = datum.value
	}

	return
}

func (c *Cache) MustGet(key string) (value interface{}) {
	value, found := c.Get(key)
	if !found {
		panic(ErrNotFound)
	}

	return
}

func (c *Cache) Has(key string) (found bool) {
	defer c.mutex.RUnlock()
	c.mutex.RLock()

	_, found = c.data[key]

	return
}

func (c *Cache) Len() int {
	defer c.mutex.RUnlock()
	c.mutex.RLock()

	return len(c.data)
}

func (c *Cache) String() string {
	defer c.mutex.RUnlock()
	c.mutex.RLock()

	var builder strings.Builder

	builder.WriteByte('[')

	var i int
	for key := range c.data {
		if i++; i > 1 {
			builder.WriteString(", ")
		}

		escapedKey := strings.ReplaceAll(key, ",", "\\,")
		builder.WriteString(escapedKey)
	}

	builder.WriteByte(']')

	return builder.String()
}
