package utils

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Retry attempts func f, in attempts, backing exponentially off by the specified Duration after each attempt which failed
// until the max number of attempts have been tried, then return an error
func Retry(attempts int, sleep time.Duration, f func() error) error {
	if err := f(); err != nil {
		if s, ok := err.(stop); ok {
			return s.error
		}

		if attempts--; attempts > 0 {
			jitter := time.Duration(rand.Int63n(int64(sleep)))
			sleep = sleep + jitter/2

			time.Sleep(sleep)
			return Retry(attempts, 2*sleep, f)
		}
		return err
	}
	return nil
}

type stop struct {
	error
}
