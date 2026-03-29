package session

import (
	"math"
	"time"
)

const (
	MaxRenews = math.MaxInt
	MaxAgeMax = 3 * 24 * time.Hour
)

func ProcessRenewed(renewed int) (int, error) {
	renewed = max(renewed+1, 0)
	if renewed > MaxRenews {
		return 0, ErrUnrenewable
	}

	return renewed, nil
}

func ProcessMaxAge(max_age time.Duration) (time.Time, error) {
	if max_age > MaxAgeMax {
		return time.Time{}, ErrTooLong
	}

	expires := time.Now().Add(max_age)
	return expires, nil
}
