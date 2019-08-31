package cache

import (
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

	expiryDuration := c.options.ExpiryDuration
	maxValues := c.options.MaxValues

	beyondMaxCount := len(c.data) - maxValues

	if expiryDuration > 0 {
		for key, datum := range c.data {
			if time.Since(datum.setTime) >= expiryDuration {
				delete(c.data, key)
				beyondMaxCount--
			}
		}
	}

	if !c.mutated {
		return
	}

	if beyondMaxCount > 0 && maxValues > 0 {
		if beyondMaxCount == 1 {
			var earliestDatum *cacheDatum

			for _, datum := range c.data {
				if earliestDatum == nil {
					earliestDatum = datum
				} else {
					earliestSetTime := earliestDatum.setTime

					if isZero := earliestSetTime.IsZero(); isZero || datum.setTime.Sub(earliestSetTime) > 0 {
						earliestDatum = datum

						if isZero {
							break
						}
					}
				}
			}

			delete(c.data, earliestDatum.key)
		} else {
			dataLen := len(c.data)
			dataMaxIndex := dataLen - 1
			sortedData := make(cacheDataSlice, dataLen)

			i := dataMaxIndex
			for _, datum := range c.data {
				sortedData[i] = datum
				i--
			}

			sort.Sort(sortedData)

			i = dataMaxIndex
			for l := i - beyondMaxCount; i > l; i-- {
				delete(c.data, sortedData[i].key)
			}
		}
	}

	c.mutated = false
}

const (
	watchSleepDurationMin = time.Millisecond
	watchSleepDurationMax = time.Minute
)

// runs in own goroutine
func (cache *Cache) watch() {
	for {
		sleepDuration := cache.options.ExpiryDuration / 10

		if sleepDuration < watchSleepDurationMin {
			sleepDuration = watchSleepDurationMin
		} else if sleepDuration > watchSleepDurationMax {
			sleepDuration = watchSleepDurationMax
		}

		for {
			time.Sleep(sleepDuration)
			cache.Clean()

			cache.mutex.RLock()
			m := cache.mutated
			cache.mutex.RUnlock()

			if m {
				break
			}
		}
	}
}
