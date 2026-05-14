package money_test

import (
	"testing"

	. "github.com/alan-b-lima/almodon/pkg/money"
)

func TestString(t *testing.T) {
	type Test struct {
		Money    Money
		Expected string
	}

	tests := []Test{
		{Money: 0, Expected: "0.00"},
		{Money: 1, Expected: "0.01"},
		{Money: 10, Expected: "0.10"},
		{Money: -10, Expected: "-0.10"},
		{Money: 1234567890, Expected: "12345678.90"},
		{Money: -1234567890, Expected: "-12345678.90"},
	}

	for _, test := range tests {
		if test.Money.String() != test.Expected {
			t.Errorf("Expected %s, got %s", test.Expected, test.Money.String())
		}
	}
}
