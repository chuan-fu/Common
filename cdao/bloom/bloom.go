package bloom

import (
	"github.com/go-redis/redis/v8"
	"github.com/spaolacci/murmur3"
)

const (
	// for detailed error rate table, see http://pages.cs.wisc.edu/~cao/papers/summary-cache/node8.html
	// maps as k in the error rate table
	maps = 14
	// constantA为位图长度/预计元素个数的常量
	// bits  = constantA  * elements
	constantA = 20
)

type Bloom struct {
	bits   uint
	bitSet BitSetProvider
}

// New create a BloomFilter, store is the backed redis, key is the key for the bloom_func filter,
// bits is how many bits will be used, maps is how many hashes for each addition.
// best practices:
// elements - means how many actual elements
// when maps = 14, formula: 0.7*(bits/maps), bits = 20*elements, the error rate is 0.000067 < 1e-4
// for detailed error rate table, see http://pages.cs.wisc.edu/~cao/papers/summary-cache/node8.html
// elements为预计元素个数，真实元素数量超过预计，错误率会急剧上升
func NewBloomFilter(store redis.Cmdable, key string, elements uint) *Bloom {
	bits := elements * constantA
	return &Bloom{
		bits:   bits,
		bitSet: newRedisBitSet(store, key, bits),
	}
}

// 添加元素
func (f *Bloom) Add(data []byte) error {
	locations := f.getLocations(data)
	if err := f.bitSet.set(locations); err != nil {
		return err
	}
	return nil
}

// 校验元素是否存在
func (f *Bloom) Exists(data []byte) (bool, error) {
	locations := f.getLocations(data)
	isSet, err := f.bitSet.check(locations)
	if err != nil {
		return false, err
	}
	if !isSet {
		return false, nil
	}

	return true, nil
}

func (f *Bloom) getLocations(data []byte) []uint {
	locations := make([]uint, maps)
	for i := uint(0); i < maps; i++ {
		hashValue := hash(append(data, byte(i)))
		locations[i] = uint(hashValue % uint64(f.bits))
	}

	return locations
}

func hash(data []byte) uint64 {
	return murmur3.Sum64(data)
}
