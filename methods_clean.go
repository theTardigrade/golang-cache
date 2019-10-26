package cache

import (
	"sort"
	"time"
)

type cacheDataSlice []*cacheDatum

func (s cacheDataSlice) Len() int           { return len(s) }
func (s cacheDataSlice) Less(i, j int) bool { return s[i].setTime.Sub(s[j].setTime) > 0 }
func (s cacheDataSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

const (
	cleanDefaultMaxItemsPerSweep = 1 << 16
)

func (c *Cache) Clean() {
	var doAnotherSweep bool

	func() {
		var maxValuesPerSweep int

		defer c.mutex.Unlock()
		c.mutex.Lock()

		preDeletionFuncExists := (c.options.PreDeletionFunc != nil)
		postDeletionFuncExists := (c.options.PostDeletionFunc != nil)

		if c.options.CleanMaxValuesPerSweep != 0 {
			maxValuesPerSweep = c.options.CleanMaxValuesPerSweep
		} else {
			maxValuesPerSweep = cleanDefaultMaxItemsPerSweep
		}

		expiryDuration := c.options.ExpiryDuration
		maxValues := c.options.MaxValues

		beyondMaxCount := len(c.data) - maxValues

		if beyondMaxCount > maxValuesPerSweep {
			beyondMaxCount = maxValuesPerSweep
			doAnotherSweep = true
		}

		if expiryDuration > 0 {
			for key, datum := range c.data {
				if beyondMaxCount == 0 {
					return
				}

				if time.Since(datum.setTime) >= expiryDuration {
					if preDeletionFuncExists {
						c.options.PreDeletionFunc(key, datum.value, datum.setTime)
					}

					delete(c.data, key)
					beyondMaxCount--

					if postDeletionFuncExists {
						c.options.PostDeletionFunc(key, datum.value, datum.setTime)
					}
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

				if preDeletionFuncExists {
					c.options.PreDeletionFunc(earliestDatum.key, earliestDatum.value, earliestDatum.setTime)
				}

				delete(c.data, earliestDatum.key)

				if postDeletionFuncExists {
					c.options.PostDeletionFunc(earliestDatum.key, earliestDatum.value, earliestDatum.setTime)
				}
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
					datum := sortedData[i]

					if preDeletionFuncExists {
						c.options.PreDeletionFunc(datum.key, datum.value, datum.setTime)
					}

					delete(c.data, datum.key)

					if postDeletionFuncExists {
						c.options.PostDeletionFunc(datum.key, datum.value, datum.setTime)
					}
				}
			}
		}

		c.mutated = false
	}()

	if doAnotherSweep {
		time.Sleep(cleanDurationGeneratedMin)
		c.Clean()
	}
}

const (
	cleanDurationGeneratedMin = time.Microsecond
	cleanDurationGeneratedMax = time.Minute
)

// runs in own goroutine
func (c *Cache) watch() {
	var prevExecutionDuration time.Duration
	startTime := time.Now()

	for {
		var cleanDuration time.Duration

		func() {
			defer c.mutex.RUnlock()
			c.mutex.RLock()

			cleanDuration = c.options.CleanDuration

			if cleanDuration < 0 {
				cleanDuration = c.options.ExpiryDuration / 10

				if cleanDuration < cleanDurationGeneratedMin {
					cleanDuration = cleanDurationGeneratedMin
				} else if cleanDuration > cleanDurationGeneratedMax {
					cleanDuration = cleanDurationGeneratedMax
				}
			}
		}()

		if cleanDuration == 0 {
			return
		}

		for {
			sleepDuration := cleanDuration - prevExecutionDuration
			prevExecutionDuration = time.Since(startTime)

			if sleepDuration > 0 {
				time.Sleep(sleepDuration)
			}

			startTime = time.Now()

			c.Clean()

			if func() bool {
				defer c.mutex.RUnlock()
				c.mutex.RLock()

				return c.mutated
			}() {
				break
			}
		}
	}
}
