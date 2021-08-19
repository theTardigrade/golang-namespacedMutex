package namespacedMutex

import (
	"sync"

	hash "github.com/theTardigrade/golang-hash"
)

func (d *Datum) masterMutexHash(namespace string) int {
	namespaceHash := hash.Uint64String(namespace)
	count := uint64(d.masterMutexesBucketCount)

	return int(namespaceHash % count)
}

func (d *Datum) masterMutex(namespace string) *sync.Mutex {
	hash := d.masterMutexHash(namespace)
	return d.masterMutexes[hash]
}
