package cache

import (
	"sort"
	"time"

	tasks "github.com/theTardigrade/golang-tasks"
)

type cacheDataSlice []*cacheDatum

func (s cacheDataSlice) Len() int           { return len(s) }
func (s cacheDataSlice) Less(i, j int) bool { return s[i].setTime.Sub(s[j].setTime) > 0 }
func (s cacheDataSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

const (
	cleanDefaultMaxItemsPerSweep = 1 << 16
)

// clean must only be called when mutex is already locked.
func (c *Cache) clean() (cleanedFully bool) {
	var maxValuesPerSweep int

	if c.options.CleanMaxValuesPerSweep > 0 {
		maxValuesPerSweep = c.options.CleanMaxValuesPerSweep
	} else {
		maxValuesPerSweep = cleanDefaultMaxItemsPerSweep
	}

	expiryDuration := c.options.ExpiryDuration
	maxValues := c.options.MaxValues

	if maxValues < 0 {
		maxValues = 0
	}

	beyondMaxCount := len(c.data) - maxValues
	var beyondMaxCountOverflow bool

	if maxValues > 0 && beyondMaxCount > maxValuesPerSweep {
		beyondMaxCount = maxValuesPerSweep
		beyondMaxCountOverflow = true
	}

	if expiryDuration > 0 {
		var handledCount int

		for _, datum := range c.data {
			if handledCount >= maxValuesPerSweep {
				return
			}

			if time.Since(datum.setTime) >= expiryDuration {
				c.unset(datum)
				handledCount++
			}
		}

		beyondMaxCount -= handledCount
	}

	if maxValues > 0 && beyondMaxCount > 0 && c.hasStatus(statusHasMutatedSinceCleanedFully) {
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

			c.unset(earliestDatum)
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
				c.unset(datum)
			}
		}
	}

	if !beyondMaxCountOverflow {
		cleanedFully = true
		c.unsetStatus(statusHasMutatedSinceCleanedFully)
	}

	return
}

func (c *Cache) Clean() (cleanedFully bool) {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	return c.clean()
}

const (
	cleanFullyInitialMaxIterations = 1 << 10
)

func (c *Cache) cleanFully() {
	var cleanedFully bool

	for i := 0; i < cleanFullyInitialMaxIterations; i++ {
		cleanedFully = func() bool {
			defer c.mutex.Unlock()
			c.mutex.Lock()

			return c.clean()
		}()

		if cleanedFully {
			break
		}
	}

	// final full clean
	if !cleanedFully {
		func() {
			defer c.mutex.Unlock()
			c.mutex.Lock()

			for {
				if cleanedFully = c.clean(); cleanedFully {
					break
				}
			}
		}()
	}
}

const (
	cleanDurationGeneratedMin = time.Microsecond
	cleanDurationGeneratedMax = time.Hour
)

// watch runs in own goroutine
func (c *Cache) cleanWatch() {
	var taskID *tasks.Identifier
	var cleanDuration time.Duration

	for {
		doneChan := make(chan struct{})

		go func(doneChan chan<- struct{}) {
			func() {
				defer c.mutex.RUnlock()
				c.mutex.RLock()

				cleanDuration = c.options.CleanDuration

				if cleanDuration <= 0 {
					cleanDuration = c.options.ExpiryDuration / 10

					if cleanDuration < cleanDurationGeneratedMin {
						cleanDuration = cleanDurationGeneratedMin
					} else if cleanDuration > cleanDurationGeneratedMax {
						cleanDuration = cleanDurationGeneratedMax
					}
				}
			}()

			if taskID == nil {
				taskID = tasks.Set(cleanDuration, true, func(id *tasks.Identifier) {
					c.cleanFully()
				})
			} else {
				taskID.ChangeInterval(cleanDuration)
			}

			doneChan <- struct{}{}
		}(doneChan)

		func(doneChan <-chan struct{}, cleanIntervalChan <-chan struct{}) {
			var d bool

			select {
			case <-doneChan:
				d = true
			case <-cleanIntervalChan:
			}

			if d {
				<-cleanIntervalChan
			} else {
				<-doneChan
			}
		}(doneChan, c.cleanIntervalChan)
	}
}

// startWatchIfNecessary must only be called when mutex is already locked.
func (c *Cache) startCleanWatchIfNecessary() {
	if !c.hasStatus(statusHasCleanWatchStarted) {
		if c.options.MaxValues > 0 || c.options.ExpiryDuration > 0 {
			go c.cleanWatch()

			c.setStatus(statusHasCleanWatchStarted)
		}
	}
}
