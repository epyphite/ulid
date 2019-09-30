package ulid

import (
	"math/rand"
	"testing"
	"time"
)

func TestUlid(t *testing.T) {

	// reproducible entropy source
	entropy := rand.New(rand.NewSource(time.Unix(1000000, 0).UnixNano()))

	// sub-ms safe ULID generator
	ulidSource := NewMonotonicULIDsource(entropy)

	for i := 0; i < 100000; i++ {
		now := time.Now()
		id, _ := ulidSource.New(now)
		println(id.String())
	}
}
