package utils

import (
	"math/rand"
	"time"
)

func AddRandomTime(duration time.Duration) (randomTime time.Duration, actual time.Duration) {
	if duration < 5*time.Minute {
		// Add a random time between 0 and 3 minutes
		rand.Seed(time.Now().UnixNano())
		randomTime = time.Duration(rand.Int63n(int64(2*time.Minute))) + 3*time.Minute
		duration += randomTime
	} else {
		// Add or subtract a random time between 0 and 3 minutes
		rand.Seed(time.Now().UnixNano())
		randomTime = time.Duration(rand.Int63n(int64(2*time.Minute))) + 3*time.Minute
		if rand.Intn(2) == 0 {
			duration += randomTime
		} else {
			duration -= randomTime
			randomTime = 0 - randomTime
		}
	}
	return randomTime, duration
}
