package session

import (
	"fmt"
	"math"
	"strconv"
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

const TokenLen = 32

type Token [TokenLen]byte

func NewToken() Token {
	var token Token
	read(token[:])

	return token
}

func (t *Token) Bytes() []byte {
	return t[:]
}

var _Format = `%0` + strconv.Itoa(2*TokenLen) + `x`

func (t Token) String() string {
	return fmt.Sprintf(_Format, t)
}

func FromString(string string) (Token, error) {
	if 2*TokenLen != len(string) {
		return Token{}, ErrInvalidToken
	}

	var token Token
	for i := range TokenLen {
		b, err := strconv.ParseInt(string[:2], 16, 8)
		if err != nil {
			return Token{}, ErrInvalidToken
		}

		token[i] = byte(b)
	}

	return token, nil
}
