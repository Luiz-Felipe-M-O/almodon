package uuid_test

import (
	"crypto/rand"
	"sync"
	"testing"
	"time"

	. "github.com/alan-b-lima/almodon/pkg/uuid"
)

func TestInvariant(t *testing.T) {
	const (
		mask_4bit = (1 << 4) - 1
	)

	const numTests = 100

	for range numTests {
		uuid := NewUUIDv7()

		var (
			unix_ts_ms = 0 |
				uint64(uuid[0x0])<<0x28 | uint64(uuid[0x1])<<0x20 |
				uint64(uuid[0x2])<<0x18 | uint64(uuid[0x3])<<0x10 |
				uint64(uuid[0x4])<<0x08 | uint64(uuid[0x5])<<0x00

			version = uint64(uuid[0x6]) >> 4
			variant = uint64(uuid[0x8]) >> 6
		)

		if version != 7 {
			t.Errorf("unexpected version, expected 7, got %d", version)
			continue
		}

		if variant != 0b10 {
			t.Errorf("unexpected version, expected 10, got %02b", variant)
			continue
		}

		t.Logf("valid UUIDv7 %v from %s", uuid, time.UnixMilli(int64(unix_ts_ms)).Format(time.RFC1123))
	}
}

func TestConcurrentUUIDGeneration(t *testing.T) {
	const numBatches, batchSize = 123, 1999
	limit := numBatches * batchSize

	result := make([]UUID, limit)
	var wg sync.WaitGroup

	wg.Add(numBatches)
	for i := range numBatches {
		offset := i * batchSize
		r := result[offset : offset+batchSize]

		go func() {
			for i := range batchSize {
				uuid := NewUUIDv7()
				r[i] = uuid
			}

			wg.Done()
		}()
	}

	wg.Wait()

	set := make(map[UUID]struct{}, limit)
	for _, v := range result {
		set[v] = struct{}{}
	}

	if len(set) < limit {
		t.Error("equal UUIDs have been generated")
	}
}

func TestInversabilityBetweenStringAndFromString(t *testing.T) {
	const numTests = 1000

	for range numTests {
		var uuid UUID
		rand.Read(uuid[:])

		str := uuid.String()
		if uuid2, err := FromString(str); err != nil {
			t.Error(err)
		} else if uuid != uuid2 {
			t.Errorf("%x and %x should be equal", uuid, uuid2)
		}
	}
}
