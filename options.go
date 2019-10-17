package cache

import "time"

type Options struct {
	ExpiryDuration         time.Duration
	MaxValues              int
	CleanDuration          time.Duration
	CleanMaxValuesPerSweep int
	PreDeletionFunc        CallbackFunc
	PostDeletionFunc       CallbackFunc
}
