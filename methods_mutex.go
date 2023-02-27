package namespacedMutex

import (
	"sync"

	hash "github.com/theTardigrade/golang-hash"
)

func (d *Datum) mutexHashFromNamespace(namespace string) int {
	keyHash := hash.Uint64String(namespace)
	count := uint64(d.mutexesBucketCount)

	return int(keyHash % count)
}

func (d *Datum) mutexFromHash(hash int) *sync.RWMutex {
	return d.mutexes[hash]
}

func (d *Datum) mutexFromNamespace(namespace string) *sync.RWMutex {
	hash := d.mutexHashFromNamespace(namespace)

	return d.mutexFromHash(hash)
}
