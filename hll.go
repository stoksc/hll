package hll

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"math"
	"math/bits"
	"sync/atomic"
)

// HyperLogLog is a specialized data structure for estimating the cardinality of a dataset.
// By default, it is divided into 256 substreams where the stream is decided by the first half byte
// of the hash.
type HyperLogLog struct {
	substreams [M]int32
}

// New creates a new HyperLogLog
func New() HyperLogLog {
	return HyperLogLog{
		substreams: [M]int32{0},
	}
}

const (
	// Am is a normalization constant used in the cardinality formula
	Am = 0.71827259325
	// M is the number of substreams
	M = 256
	// HCC is the point at which high cardinality corrections begin
	HCC = 143165576
)

// Cardinality estimated the number of unique elements inserted into the HyperLogLog
func (h *HyperLogLog) Cardinality() uint64 {
	s := 0.0
	for j := range h.substreams {
		s += math.Pow(2, float64(-h.substreams[j]))
	}

	s = math.Pow(s, -1)

	e := Am * M * M * s

	var es float64

	if e <= 2.5*M {
		cz := h.substreamsAtZero()
		if cz != 0 {
			es = h.linearCounting(cz)
		} else {
			es = e
		}
	} else if e <= HCC {
		es = e
	} else {
		es = -4294967296 * math.Log(1-(e/4294967296))
	}

	return uint64(es)
}

func (h *HyperLogLog) substreamsAtZero() int {
	cz := 0

	for _, ss := range h.substreams {
		if ss == 0 {
			cz++
		}
	}

	return cz
}

func (h *HyperLogLog) linearCounting(cz int) float64 {
	return M * math.Log(M/float64(cz))
}

// Insert inserts a new element into the HyperLogLog
func (h *HyperLogLog) Insert(x interface{}) {
	b := []byte(fmt.Sprintf("%v", x))
	hash := md5.Sum(b)

	head := hash[0]
	tail := hash[1:]

	newh := binary.BigEndian.Uint64(tail)

	newlz := int32(bits.LeadingZeros64(newh)) + 1

casf:
	curlz := h.substreams[head]
	if curlz < newlz && !atomic.CompareAndSwapInt32(&h.substreams[head], curlz, newlz) {
		goto casf
	}
}
