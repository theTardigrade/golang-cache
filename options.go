package cache

import "time"

type Options struct {
	ExpiryDuration time.Duration
	MaxValues      int
}
