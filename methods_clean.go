package cache

import (
	"math/rand"
	"sort"
	"time"
)

type cacheDataSlice []*cacheDatum

func (s cacheDataSlice) Len() int           { return len(s) }
func (s cacheDataSlice) Less(i, j int) bool { return s[i].setTime.Sub(s[j].setTime) > 0 }
func (s cacheDataSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func (c *Cache) Clean() {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	beyondMaxCount := len(c.data) - c.maxValues

	if c.expiryDuration >= 0 {
		for key, datum := range c.data {
			if time.Since(datum.setTime) >= c.expiryDuration {
				delete(c.data, key)
				beyondMaxCount--
			}
		}
	}

	if c.maxValues >= 0 && beyondMaxCount > 0 {
		if beyondMaxCount == 1 {
			var earliestDatum *cacheDatum

			for _, datum := range c.data {
				earliestSetTime := earliestDatum.setTime
				if isZero := earliestSetTime.IsZero(); isZero || datum.setTime.Sub(earliestSetTime) > 0 {
					earliestDatum = datum

					if isZero {
						break
					}
				}
			}

			delete(c.data, earliestDatum.key)
		} else {
			dataLen := len(c.data)
			sortedData := make(cacheDataSlice, 0, dataLen)

			for _, datum := range c.data {
				sortedData = append(sortedData, datum)
			}

			sort.Sort(sortedData)

			i := dataLen - 1
			l := i - beyondMaxCount
			for i > l {
				delete(c.data, sortedData[i].key)
				i--
			}
		}
	}
}

// runs in own goroutine
func (cache *Cache) watch() {
	sleepDuration := cache.expiryDuration / 6
	if sleepDuration < time.Minute {
		sleepDuration = time.Minute
	} else if sleepDuration > time.Hour {
		sleepDuration = time.Hour
	}

	initialSleepDuration := time.Duration(rand.Intn(int(sleepDuration/time.Millisecond))) * time.Millisecond

	for time.Sleep(initialSleepDuration); ; time.Sleep(sleepDuration) {
		cache.Clean()
	}
}
