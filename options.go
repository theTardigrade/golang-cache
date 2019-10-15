package cache

import "time"

type Options struct {
	ExpiryDuration  time.Duration
	MaxValues       int
	CleanDuration   time.Duration
	PreDeletionFunc CallbackFunc
}
