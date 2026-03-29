package resource

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
)

var (
	errNotPointer  = errors.New("query: v must be a pointer type")
	errNotToStruct = errors.New("query: v must be a pointer type")
	errUnsettable  = errors.New("query: cannot change contents of v")
)

func QueryParams(q url.Values, v any) error {
	rt := reflect.ValueOf(v)

	if rt.Kind() != reflect.Pointer {
		return errNotPointer
	}

	if rt.Elem().Kind() != reflect.Struct {
		return errNotToStruct
	}

	return queryParams(q, v)
}

func queryParams(q url.Values, v any) error {
	rt := reflect.TypeOf(v).Elem()
	rv := reflect.ValueOf(v).Elem()

	if !rv.CanSet() {
		return errUnsettable
	}

	for i := range rt.NumField() {
		field := rt.Field(i)
		query := field.Tag.Get("query")

		if !q.Has(query) {
			continue
		}

		val := q.Get(query)

		switch t := field.Type; t.Kind() {
		case reflect.String:
			rv.Field(i).SetString(val)

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			num, err := strconv.ParseInt(val, 10, int(t.Size())*8)
			if err != nil {
				return fmt.Errorf("query: not convertible to an %v: %w", t, err)
			}

			rv.Field(i).SetInt(int64(num))

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			num, err := strconv.ParseUint(val, 10, int(t.Size())*8)
			if err != nil {
				return fmt.Errorf("query: not convertible to an %v: %w", t, err)
			}

			rv.Field(i).SetUint(uint64(num))
		}
	}

	return nil
}
