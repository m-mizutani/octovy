package utils

import (
	"math"
	"time"

	"github.com/m-mizutani/goerr"
)

var ErrBackoffLimitExceeded = goerr.New("Backoff limit is exceeded")
var BackoffBaseWaitTime float64 = 1.25

// Backoff provides exponential backoff. First return value of task means 'isExit' and leave backoff loop if true. If loop count reaches to limit, Backoff() returns ErrBackoffLimitExceeded.
func Backoff(limit int, task func() (bool, error)) error {
	for i := 0; i < limit; i++ {
		isExit, err := task()
		if err != nil {
			return err
		}
		if isExit {
			return nil
		}
		wait := math.Pow(BackoffBaseWaitTime, float64(i))
		time.Sleep(time.Millisecond * time.Duration(wait*1000))
	}

	return goerr.Wrap(ErrBackoffLimitExceeded).With("limit", limit)
}
