package helpers

import (
	"time"
)

func DurationInHours(i int) time.Duration {
	return time.Duration(i) * time.Hour
}
