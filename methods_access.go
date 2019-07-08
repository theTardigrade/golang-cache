package cache

import "strings"

func (c *Cache) Get(key string) (interface{}, bool) {
	defer c.mutex.RUnlock()
	c.mutex.RLock()

	datum, ok := c.data[key]
	if !ok {
		return nil, false
	}

	return datum.value, true
}

func (c *Cache) MustGet(key string) interface{} {
	value, ok := c.Get(key)
	if !ok {
		panic(ErrNotFound)
	}

	return value
}

func (c *Cache) Has(key string) bool {
	defer c.mutex.RUnlock()
	c.mutex.RLock()

	_, ok := c.data[key]

	return ok
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

		builder.WriteString(key)
	}

	builder.WriteByte(']')

	return builder.String()
}
