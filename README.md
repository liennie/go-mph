# go-mph

minimal perfect hash functions

godoc: https://pkg.go.dev/github.com/liennie/go-mph

## Changes

* some memory optimizations, making `New` run faster due to less allocations
* using a second hash instead of `xorshiftmult`, making `New` more robust against collisions,
but also making it slower and making `Query` slightly slower
* sanity checks and duplicate key check

## Benchmarks

Four versions were benchmarked: without change, with memory optimizations,
with a second hash instead of xorshiftmult and with duplicate key detection.
Two key sets were used: small with 26 keys and big with 1 094 068 keys.

```
goos: linux
goarch: amd64
pkg: github.com/liennie/go-mph
cpu: 13th Gen Intel(R) Core(TM) i7-1365U
```

### New (small)

```
                        │ sec/op       vs base         │
original.new.small.txt    3.590µ ± 2%
optimized.new.small.txt   2.745µ ± 2%  -23.55% (n=100)
hash.new.small.txt        2.626µ ± 2%  -26.85% (n=100)
dedup.new.small.txt       2.738µ ± 2%  -23.73% (n=100)

                        │ B/op          vs base         │
original.new.small.txt    2.508Ki ± 0%
optimized.new.small.txt   2.039Ki ± 0%  -18.69% (n=100)
hash.new.small.txt        2.039Ki ± 0%  -18.69% (n=100)
dedup.new.small.txt       2.039Ki ± 0%  -18.69% (n=100)

                        │ allocs/op   vs base         │
original.new.small.txt    47.00 ± 0%
optimized.new.small.txt   19.00 ± 0%  -59.57% (n=100)
hash.new.small.txt        19.00 ± 0%  -59.57% (n=100)
dedup.new.small.txt       19.00 ± 0%  -59.57% (n=100)
```

### New (big)

```
                      │ sec/op       vs base         │
original.new.big.txt    300.0m ± 1%
optimized.new.big.txt   201.3m ± 1%  -32.92% (n=100)
hash.new.big.txt        284.7m ± 1%  -5.10% (n=100)
dedup.new.big.txt       299.6m ± 1%  ~ (p=0.761 n=100)

                      │ B/op          vs base         │
original.new.big.txt    172.8Mi ± 0%
optimized.new.big.txt   93.66Mi ± 0%  -45.79% (n=100)
hash.new.big.txt        93.66Mi ± 0%  -45.79% (n=100)
dedup.new.big.txt       93.66Mi ± 0%  -45.79% (n=100)

                      │ allocs/op    vs base         │
original.new.big.txt    1.496M ± 0%
optimized.new.big.txt   440.6k ± 0%  -70.55% (n=100)
hash.new.big.txt        440.6k ± 0%  -70.55% (n=100)
dedup.new.big.txt       440.6k ± 0%  -70.55% (n=100)
```

### Query (small)

```
                          │ sec/op       vs base                 │
original.query.small.txt    125.5n ± 0%
optimized.query.small.txt   125.3n ± 0%  ~ (p=0.673 n=100)
hash.query.small.txt        156.0n ± 0%  +24.25% (p=0.000 n=100)
dedup.query.small.txt       156.4n ± 0%  +24.61% (p=0.000 n=100)
```

### Query (big)

```
                        │ sec/op       vs base                │
original.query.big.txt    69.68m ± 1%
optimized.query.big.txt   68.25m ± 1%  -2.05% (p=0.000 n=100)
hash.query.big.txt        74.41m ± 1%  +6.79% (p=0.000 n=100)
dedup.query.big.txt       72.79m ± 0%  +4.47% (p=0.000 n=100)
```
