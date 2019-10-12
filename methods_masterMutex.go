package namespacedMutex

import (
	"fmt"
	"sync"
)

func (d *Datum) masterMutexHash(namespace string) int {
	var hash uint = 0x811c9dc5
	count := uint(d.masterMutexesBucketCount)

	for i := len(namespace) - 1; i >= 0; i-- {
		b := namespace[i]
		hash ^= uint(b)
		hash *= 0x01000193
	}

	fmt.Println("DEBUG: MASTER MUTEX", namespace, hash%count, hash)

	return int(hash % count)
}

func (d *Datum) masterMutex(namespace string) *sync.Mutex {
	hash := d.masterMutexHash(namespace)
	return d.masterMutexes[hash]
}
