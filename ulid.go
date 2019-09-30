package ulid

import (
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/oklog/ulid"
)

type MonotonicULIDsource struct {
	sync.Mutex
	entropy  *rand.Rand
	lastMs   uint64
	lastULID ulid.ULID
}

func NewMonotonicULIDsource(entropy *rand.Rand) *MonotonicULIDsource {
	initial, err := ulid.New(ulid.Now(), entropy)
	if err != nil {
		log.Println("Error Generating a new source")
		return nil
	}

	return &MonotonicULIDsource{
		entropy:  entropy,
		lastMs:   0,
		lastULID: initial,
	}
}

func (u *MonotonicULIDsource) New(t time.Time) (ulid.ULID, error) {
	u.Lock()
	defer u.Unlock()

	ms := ulid.Timestamp(t)
	var err error
	if ms > u.lastMs {
		u.lastMs = ms
		u.lastULID, err = ulid.New(ms, u.entropy)
		return u.lastULID, err
	}
	// if the ms are the same, increment the entropy part of the last ULID
	// and use it as the entropy part of the new ULID.
	incrEntropy := incrementBytes(u.lastULID.Entropy())
	var dup ulid.ULID
	dup.SetTime(ms)
	if err := dup.SetEntropy(incrEntropy); err != nil {
		return dup, err
	}
	u.lastULID = dup
	u.lastMs = ms
	return dup, nil
}
func incrementBytes(in []byte) []byte {
	const (
		minByte byte = 0
		maxByte byte = 255
	)

	// copy the byte slice
	out := make([]byte, len(in))
	copy(out, in)

	// iterate over the byte slice from right to left
	// most significant byte == first byte (big-endian)
	leastSigByteIdx := len(out) - 1
	mostSigByteIdex := 0
	for i := leastSigByteIdx; i >= mostSigByteIdex; i-- {

		// If its value is 255, rollover back to 0 and try the next byte.
		if out[i] == maxByte {
			out[i] = minByte
			continue
		}
		// Else: increment.
		out[i]++
		return out
	}
	// panic if the increments are exhausted

	log.Println("Panic Resources for ms exhausted")
	return nil
}
