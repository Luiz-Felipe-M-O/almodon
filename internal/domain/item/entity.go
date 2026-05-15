package item

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/alan-b-lima/almodon/pkg/money"
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

func ProcessUnitCost(cost money.Money) (money.Money, error) {
	if cost < 0 {
		return 0, ErrUnitCostNegative
	}
	return cost, nil
}

func ProcessExpires(expires time.Time) (time.Time, error) {
	return expires, nil
}

func StatusAmount(amount, min float64) Stock {
	switch {
	case min <= 0:
		return StockFine
	case amount <= min:
		return StockWarning
	}

	return StockFine
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

var stocks = [...]string{
	StockFine:    "fine",
	StockWarning: "warning",
	StockEmpty:   "empty",
}

func (s Stock) String() string {
	if 0 <= int(s) && int(s) < len(stocks) {
		str := stocks[s]
		if str != "" {
			return str
		}
	}

	return "stock(" + strconv.Itoa(int(s)) + ")"
}

func (s Stock) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

type Expiration int

const (
	ExpirationNone Expiration = iota
	ExpirationFine
	ExpirationWarning
	ExpirationExpired
)

var expirations = [...]string{
	ExpirationNone:    "none",
	ExpirationFine:    "fine",
	ExpirationWarning: "warning",
	ExpirationExpired: "expired",
}

func (s Expiration) String() string {
	if 0 <= int(s) && int(s) < len(expirations) {
		str := expirations[s]
		if str != "" {
			return str
		}
	}

	return "expiration(" + strconv.Itoa(int(s)) + ")"
}

func (s Expiration) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}
