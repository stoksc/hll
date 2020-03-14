package pds

import (
	"fmt"
	"hash"
	"hash/fnv"
)

// BloomFilter is a probabilistic data structure, used to
// test set membership (and also cardinatlity estimation).
type BloomFilter struct {
	buf [m / 4]byte
}

const (
	// M is the size of the bloom filter
	m = 256
	// K is the number of hash functions
	K = 3
)

var hashers = []hash.Hash{
	fnv.New32(),
	fnv.New64(),
	fnv.New128(),
}

// NewBloomFilter creates a new bloom filter
func NewBloomFilter() BloomFilter {
	return BloomFilter{
		buf: [m / 4]byte{},
	}
}

// Insert inserts a value into the bloom filter
func (b *BloomFilter) Insert(x interface{}) {
	bs := []byte(fmt.Sprintf("%v", x))

	for _, h := range hashers {
		sum := h.Sum(bs)
		isum := int(sum[0])
		inbf := isum % m

		bucket := inbf / 4
		bit := inbf % 4

		switch bit {
		case 0:
			b.buf[bucket] = b.buf[bucket] | 1
		case 1:
			b.buf[bucket] = b.buf[bucket] | 2
		case 2:
			b.buf[bucket] = b.buf[bucket] | 4
		case 3:
			b.buf[bucket] = b.buf[bucket] | 8
		}
	}
}

// Check checks if a vale is in the bloom filter
func (b *BloomFilter) Check(x interface{}) bool {
	bs := []byte(fmt.Sprintf("%v", x))

	isIn := true
	for _, h := range hashers {
		sum := h.Sum(bs)
		isum := int(sum[0])
		inbf := isum % m

		bucket := inbf / 4
		bit := inbf % 4

		switch bit {
		case 0:
			isIn = isIn && (b.buf[bucket]&1 == 1)
		case 1:
			isIn = isIn && (b.buf[bucket]&2 == 2)
		case 2:
			isIn = isIn && (b.buf[bucket]&4 == 4)
		case 3:
			isIn = isIn && (b.buf[bucket]&8 == 8)
		}
	}

	return isIn
}
