package sse_test

import (
	"bytes"
	"math/rand/v2"
	"net/http"
	"slices"
	"testing"

	. "github.com/alan-b-lima/almodon/pkg/sse"
)

type MockResponseWriter struct {
	buf []byte
}

func (m *MockResponseWriter) Header() http.Header { return http.Header{} }
func (m *MockResponseWriter) WriteHeader(int)     {}
func (m *MockResponseWriter) Flush()              {}

func (m *MockResponseWriter) Write(b []byte) (int, error) {
	m.buf = append(m.buf, b...)
	return len(b), nil
}

func TestWrite(t *testing.T) {
	type Template struct {
		write   []byte
		expect  []byte
		samples int
	}

	type Test struct {
		writes [][]byte
		expect []byte
	}

	rand := rand.New(rand.NewPCG(0, 2))

	templates := []Template{
		{
			write:   []byte("hello"),
			expect:  []byte("data: hello\r\n\r\n"),
			samples: 30,
		},
		{
			write:   []byte("hello\nworld"),
			expect:  []byte("data: hello\r\ndata: world\r\n\r\n"),
			samples: 30,
		},
		{
			write:   []byte("hello\nworld\r\neverybody!"),
			expect:  []byte("data: hello\r\ndata: world\r\ndata: everybody!\r\n\r\n"),
			samples: 30,
		},
		{
			write:   []byte("hello\n\rwo\rrld\r\neverybody!"),
			expect:  []byte("data: hello\r\ndata: world\r\ndata: everybody!\r\n\r\n"),
			samples: 30,
		},
	}

	var tests []Test
	for _, template := range templates {
		for range template.samples {
			tests = append(tests, Test{
				writes: fragment(template.write, 1+rand.IntN(len(template.write)), rand),
				expect: template.expect,
			})
		}
	}

	for _, test := range tests {
		m := &MockResponseWriter{}
		sse, err := New(m)
		if err != nil {
			t.Errorf("New should not error: %v", err)
		}

		for _, write := range test.writes {
			_, err := sse.Write(write)
			if err != nil {
				t.Errorf("Write should not error: %v", err)
			}
		}

		_, err = sse.Dispatch()
		if err != nil {
			t.Errorf("Dispatch should not error: %v", err)
		}

		if !bytes.Equal(m.buf, test.expect) {
			t.Errorf(
				"unexpected output:\n\twrites:   %+q\n\tobtained: %+q\n\texpected: %+q",
				test.writes, m.buf, test.expect,
			)
		}
	}
}

func fragment(expr []byte, pieces int, rand *rand.Rand) [][]byte {
	if len(expr) < pieces {
		pieces = len(expr)
	}

	cuts := make([]int, 1, pieces+1)
	for range pieces - 1 {
		cuts = append(cuts, rand.IntN(len(expr)))
	}
	cuts = append(cuts, len(expr))

	slices.Sort(cuts)

	fragments := make([][]byte, 0, pieces)
	for i := range pieces {
		lo, hi := cuts[i], cuts[i+1]
		fragments = append(fragments, expr[lo:hi])
	}

	return fragments
}
