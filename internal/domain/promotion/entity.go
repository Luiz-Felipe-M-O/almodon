package promotion

import "time"

const MaxAgeMax = 3 * 24 * time.Hour

func ProcessMaxAge(max_age time.Duration) (time.Time, error) {
	if max_age > MaxAgeMax {
		return time.Time{}, ErrTooLong
	}

	return time.Now().Add(max_age), nil
}
