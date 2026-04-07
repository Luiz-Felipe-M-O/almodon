package query_test

import (
	"errors"
	"net/url"
	"reflect"
	"testing"
	"time"

	. "github.com/alan-b-lima/almodon/pkg/query"
)

func TestRejectsInvalidInput(t *testing.T) {
	type Test struct {
		name     string
		input    any
		expected error
	}

	tests := []Test{
		{name: "nil-pointer", input: nil, expected: ErrNilPointer},
		{name: "non-pointer", input: struct{}{}, expected: ErrNotPointerToStruct},
		{name: "pointer to non-struct", input: new(int), expected: ErrNotPointerToStruct},
		{name: "nil pointer to struct", input: (*struct{ Name string })(nil), expected: ErrNotPointerToStruct},
	}

	for _, test := range tests {
		err := QueryParams(url.Values{"name": {"alice"}}, test.input)
		if err == nil {
			t.Errorf("query test %+q shouldn't have succeded", test.name)
			continue
		}

		if !errors.Is(err, test.expected) {
			t.Errorf("%+q: %v", test.name, err)
		}
	}
}

func TestParsing(t *testing.T) {
	type query_input struct {
		Name    string    `query:"name"`
		Age     int       `query:"age"`
		Balance uint32    `query:"balance"`
		Hidden  bool      `query:"hidden"`
		Tags    []string  `query:"tag"`
		Numbers []int     `query:"num"`
		Created time.Time `query:"created"`
	}

	type Test struct {
		given    url.Values
		expected query_input
	}

	timestamp := time.Date(2026, time.March, 14, 15, 92, 65, 0, time.UTC)

	tests := []Test{
		{
			given: url.Values{
				"name":    {"Luan"},
				"age":     {"23"},
				"balance": {"900"},
				"hidden":  {"false"},
				"tag":     {"admin", "ops"},
				"num":     {"7", "11", "13"},
				"created": {timestamp.Format(time.RFC3339Nano)},
			},
			expected: query_input{
				Name:    "Luan",
				Age:     23,
				Balance: 900,
				Hidden:  false,
				Tags:    []string{"admin", "ops"},
				Numbers: []int{7, 11, 13},
				Created: timestamp,
			},
		},
		{
			given: url.Values{
				"name":    {"Mateus"},
				"age":     {"30"},
				"balance": {"65"},
				"hidden":  {"T"},
				"tag":     {"admin", "ops"},
				"created": {timestamp.Format(time.RFC3339)},
			},
			expected: query_input{
				Name:    "Mateus",
				Age:     30,
				Balance: 65,
				Hidden:  true,
				Tags:    []string{"admin", "ops"},
				Created: timestamp,
			},
		},
	}

	for _, test := range tests {
		var v query_input
		if err := QueryParams(test.given, &v); err != nil {
			t.Error(err)
			continue
		}

		if !reflect.DeepEqual(v, test.expected) {
			t.Errorf("Unexpected result:\n\tgot:  %v\n\twant: %v", v, test.expected)
		}
	}
}
