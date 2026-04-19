package item

import (
	"encoding/json"
	"time"
)

const (
	ExpiresWarnThreshold = 30 * 24 * time.Hour
)

func ProcessAmount(amount float64) (float64, error) {
	if amount < 0 {
		return 0, ErrUnitCostNegative
	}
	return amount, nil
}

func StatusAmount(amount , min float64) Stock {
	switch {
	case min <= 0:
		return StockFine
	case amount <= min:
		return StockWarning
	}

	return StockFine
}

func ProcessUnitCost(cost float64) (float64, error) {
	if cost < 0 {
		return 0, ErrUnitCostNegative
	}
	return cost, nil
}

func ProcessArrival(arrival time.Time) (time.Time, error) {
	return arrival, nil
}

func ProcessExpires(expires time.Time) (time.Time, error) {
	return expires, nil
}

func StatusExpires(expires time.Time) Expiration {
	if expires.IsZero() {
		return ExpirationNone
	}

	switch diff := time.Until(expires); {
	case diff <= 0:
		return ExpirationExpired
	case diff < ExpiresWarnThreshold:
		return ExpirationWarning
	default:
		return ExpirationFine
	}
}

type Stock int

const (
	StockEmpty Stock = iota
	StockWarning
	StockFine
)

func (s Stock) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s Stock) String() string {
	switch s {
	case StockFine:
		return "fine"
	case StockWarning:
		return "warning"
	case StockEmpty:
		return "empty"
	}

	return ""
}

type Expiration int

const (
	ExpirationNone Expiration = iota
	ExpirationFine
	ExpirationWarning
	ExpirationExpired
)

func (s Expiration) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s Expiration) String() string {
	switch s {
	case ExpirationFine:
		return "fine"
	case ExpirationWarning:
		return "warning"
	case ExpirationExpired:
		return "expired"
	}

	return "none"
}
