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

func (d *Datum) mutexIndexFromNamespace(namespace string) int {
	indexBig := hash.Uint256String(namespace)

	return int(indexBig.Mod(indexBig, d.bucketCountBig).Uint64())
}

func (d *Datum) mutexFromIndex(index int) *sync.RWMutex {
	return d.mutexes[index]
}

func (d *Datum) mutexFromNamespace(namespace string) *sync.RWMutex {
	index := d.mutexIndexFromNamespace(namespace)

	return d.mutexes[index]
}
