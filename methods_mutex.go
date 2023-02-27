package namespacedMutex

import (
	"sync"

	hash "github.com/theTardigrade/golang-hash"
)

func (d *Datum) initMutexes() {
	bc := d.bucketCount

	d.mutexes = make([]*sync.RWMutex, bc)

	for bc--; bc >= 0; bc-- {
		d.mutexes[bc] = new(sync.RWMutex)
	}
}

func (d *Datum) mutexHashFromNamespace(namespace string) int {
	keyHash := hash.Uint256String(namespace)

	return int(keyHash.Mod(keyHash, d.bucketCountBig).Uint64())
}

func (d *Datum) mutexFromHash(hash int) *sync.RWMutex {
	return d.mutexes[hash]
}

func (d *Datum) mutexFromNamespace(namespace string) *sync.RWMutex {
	hash := d.mutexHashFromNamespace(namespace)

	return d.mutexes[hash]
}
