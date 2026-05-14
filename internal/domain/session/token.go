package session

import (
	"crypto/rand"
	"sync"
)

var (
	pool   [32 * 256]byte
	offset = len(pool)

	mu sync.Mutex
)

func read(p []byte) (int, error) {
	mu.Lock()
	defer mu.Unlock()

	if len(p) > len(pool) {
		return rand.Read(p)
	}

	if len(pool[offset:]) < len(p) {
		n := copy(p, pool[offset:])
		rand.Read(pool[n:])
		offset = 0
	}

	n := copy(p, pool[offset:])
	offset += n

	return n, nil
}
