pds is some probabilistic data structures i wrote so i remember them.

`pds/hll`
simple, correct, concurrent, lockless hyperloglog: e.g.

```go
hll := hll.New()

for i := 0; i <= 10; i++ {
    for j := 0; j < 1000; j++ {
        hll.Insert(j)
    }
}

c := hll.Cardinality()
```

`pds/bf`

```go
bf := NewBloomFilter()

bf.Insert(10)

// res = true
res := bf.Check(10)

// res probably = false
res = bf.Check(11)
```

