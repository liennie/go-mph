// Package mph implements a hash/displace minimal perfect hash function.
package mph

import (
	"math"
	"sort"

	"github.com/dgryski/go-metro"
)

// Table stores the values for the hash function
type Table struct {
	Values []uint32
	Seeds  []int32
}

type entry struct {
	idx  uint32
	hash uint64
}

// New constructs a minimal perfect hash function for the set of keys which returns the index of item in the keys array.
func New(keys []string) *Table {
	assert(len(keys) <= math.MaxInt32, "too many keys")

	size := uint64(nextPower2(len(keys)))

	pool := make([]entry, len(keys))
	h := make([][]entry, size)
	for idx, k := range keys {
		hash := metro.Hash64Str(k, 0)
		i := hash & (size - 1)
		if h[i] == nil {
			h[i] = pool[:0:1]
			pool = pool[1:]
		}
		h[i] = append(h[i], entry{uint32(idx), hash})
	}

	sort.Slice(h, func(i, j int) bool { return len(h[i]) > len(h[j]) })

	values := make([]uint32, size)
	seeds := make([]int32, size)

	assigned := 0
	entries := make(map[uint64]uint32, len(h[0]))

	var hidx int
	for hidx = 0; hidx < len(h) && len(h[hidx]) > 1; hidx++ {
		subkeys := h[hidx]

		var seed uint64
		clear(entries)

	newseed:
		for seed = 1; seed <= math.MaxInt32; seed++ {
			for _, k := range subkeys {
				hash := metro.Hash64Str(keys[k.idx], seed)
				i := hash & (size - 1)
				if entries[i] == 0 && values[i] == 0 {
					// looks free, claim it
					// idx+1 so we can identify empty entries in the table with 0
					entries[i] = k.idx + 1
					continue
				}

				// found a collision, reset and try a new seed
				clear(entries)
				continue newseed
			}

			// made it through; everything got placed
			break
		}
		assert(seed <= math.MaxInt32, "no viable seed found")

		// mark subkey spaces as claimed ...
		for k, v := range entries {
			values[k] = v
		}
		assigned += len(entries)

		// ... and assign this seed value for every subkey
		// NOTE(dgryski): While k.hash is different for each entry, i = k.hash % size is the same.
		// We don't need to loop over the entire slice, we can just take the seed from the first entry.

		i := subkeys[0].hash & (size - 1)
		seeds[i] = int32(seed)
	}

	// find the unassigned entries in the table
	h = h[hidx:]
	for dst := range values {
		if values[dst] == 0 {
			if len(h) > 0 && len(h[0]) > 0 {
				k := h[0][0]
				h = h[1:]
				i := k.hash & (size - 1)
				values[dst] = k.idx
				seeds[i] = -int32(dst + 1)
			}

		} else {
			// decrement idx as this is now the final value for the table
			values[dst]--
		}
	}

	return &Table{
		Values: values,
		Seeds:  seeds,
	}
}

// Query looks up an entry in the table and return the index.
func (t *Table) Query(k string) uint32 {
	size := uint64(len(t.Values))
	hash := metro.Hash64Str(k, 0)
	i := hash & (size - 1)
	seed := t.Seeds[i]
	if seed < 0 {
		return t.Values[-seed-1]
	}

	hash = metro.Hash64Str(k, uint64(seed))
	i = hash & (size - 1)
	return t.Values[i]
}

func nextPower2(n int) int {
	i := 1
	for i < n {
		i *= 2
	}
	return i
}

func assert(cond bool, err string) {
	if !cond {
		panic(err)
	}
}
