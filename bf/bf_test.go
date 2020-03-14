package pds

import "testing"

func TestBloomFilter(t *testing.T) {
	bf := NewBloomFilter()

	inVals := []int{10, 20, 40}

	for v := range inVals {
		bf.Insert(v)
	}

	for v := range inVals {
		if !bf.Check(v) {
			t.Errorf("bf.Check(%v) = false; expected true", v)
		}
	}

	outVals := []int{11, 21, 41}

	for v := range outVals {
		if !bf.Check(v) {
			t.Errorf("bf.Check(%v) = true; expected false", v)
		}
	}
}
