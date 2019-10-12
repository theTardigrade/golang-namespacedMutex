package namespacedMutex

import (
	"sync"
)

const (
	masterMutexHashPrime  uint64 = 0x00000100000001b3
	masterMutexHashOffset uint64 = 0xcbf29ce484222325
)

func (d *Datum) masterMutexHash(namespace string) int {
	hash := masterMutexHashOffset
	count := uint64(d.masterMutexesBucketCount)

	for i := len(namespace) - 1; i >= 0; i-- {
		b := namespace[i]
		hash ^= uint64(b)
		hash *= masterMutexHashPrime
	}

	return int(hash % count)
}

func (d *Datum) masterMutex(namespace string) *sync.Mutex {
	hash := d.masterMutexHash(namespace)
	return d.masterMutexes[hash]
}
