package user_test

import (
	"strings"
	"testing"

	. "github.com/alan-b-lima/almodon/internal/domain/user"
)

func TestProcessEmail(t *testing.T) {
	type Test struct {
		input    string
		succeeds bool
	}

	tests := []Test{
		{"usuario@almodon.com", true},
		{"nome.sobrenome@dominio.br", true},
		{"   nome.sobrenome@dominio.br", true},
		{"nome.sobrenome@dominio.br   ", true},
		{"", false},
		{"usuarioalmodon.com", false},
		{"usuario@", false},
		{"@almodon.com", false},
		{"usuario@.com  ", false},
		{"usuario@com", false},
		{"usuario@dominio.c", false},
	}

	for _, test := range tests {
		_, err := ProcessEmail(test.input)

		if (err == nil) != test.succeeds {
			if test.succeeds {
				t.Errorf("Email %+q: did not expect error, but got %+q", test.input, err)
			} else {
				t.Errorf("Email %+q: expected error, but got nil", test.input)
			}
		}
	}
}

func TestProcessPassword(t *testing.T) {
	type Test struct {
		input   string
		failure error
	}

	tests := []Test{
		{"SenhaForte123!", nil},
		{"12345678", nil},
		{"1234567", ErrPasswordTooShort},
		{strings.Repeat("a", 65), ErrPasswordTooLong},
		{" 12345678", ErrPasswordLeadOrTrailWhitespace},
		{"12345678 ", ErrPasswordLeadOrTrailWhitespace},
		{"", ErrPasswordTooShort},
		{"Senha\000123", ErrPasswordIllegalCharacters},
	}

	for _, test := range tests {
		hash, err := ProcessPassword(test.input)
		switch {
		case test.failure == nil && err == nil:
			if err := ComparePassword(hash, test.input); err != nil {
				t.Errorf("Password %+q: comparison did not work, got %+q", test.input, err)
			}

		case test.failure == nil && err != nil:
			t.Errorf("Password %+q: did not expect error, but got %+q", test.input, err)

		case test.failure != nil && err == nil:
			t.Errorf("Password %+q: expected error %+q, but got nil", test.input, test.failure)

		case test.failure != nil && err != nil:
			if test.failure != err {
				t.Errorf("Password %+q: expected error %+q, but got %+q", test.input, test.failure, err)
			}
		}
	}
}
