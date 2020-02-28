package hll

import (
	"sync"
	"testing"
)

type TestObject struct {
	int
	string
}

func TestHyperLogLog(t *testing.T) {
	testCases := []struct {
		dupe int
		uniq uint64
		err  float64
		want uint64
	}{
		{0, 10, 0.1, 10},
		{0, 100, 0.1, 100},
		{0, 1000, 0.1, 1000},
		{0, 2500, 0.1, 2500},
		{1, 2500, 0.1, 2500},
		{0, 10000, 0.1, 10000},
		{2, 10000, 0.1, 10000},
		{5, 250000, 0.1, 250000},
		{0, 7777777, 0.1, 7777777},
	}

	for _, tc := range testCases {
		hll := New()

		for i := 0; i <= tc.dupe; i++ {
			for i := uint64(0); i < tc.uniq; i++ {
				hll.Insert(i)
			}
		}

		expected, actual := tc.want, hll.Cardinality()
		err := 1 - (float64(actual) / float64(expected))
		if err > tc.err {
			t.Errorf("hll.Cardinality() = %v; want %v with an err of %v; got %v", actual, expected, tc.err, err)
		}
	}
}

func TestHyperLogLogConcurrent(t *testing.T) {
	hll := New()

	var wg sync.WaitGroup
	for i := 0; i <= 15; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := uint64(0); i < 7777777; i++ {
				hll.Insert(i)
			}
		}()
	}

	wg.Wait()

	expected, actual := 7777777, hll.Cardinality()
	err := 1 - (float64(actual) / float64(expected))
	if err > 0.1 {
		t.Errorf("hll.Cardinality() = %v; want %v with an err of %v; got %v", actual, expected, 0.1, err)
	}
}
