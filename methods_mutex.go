package namespacedMutex

import (
	"strings"
	"sync"

	hash "github.com/theTardigrade/golang-hash"
)

func (d *Datum) mutexKeyFromNamespaces(namespaces []string) string {
	switch len(namespaces) {
	case 0:
		return ""
	case 1:
		return namespaces[0]
	}

	var builder strings.Builder
	separator := d.namespaceSeparator

	for i, n := range namespaces {
		if i > 0 {
			builder.WriteString(separator)
		}

		builder.WriteString(n)
	}

	return builder.String()
}

func (d *Datum) mutexHashFromKey(key string) int {
	keyHash := hash.Uint64String(key)
	count := uint64(d.mutexesBucketCount)

	return int(keyHash % count)
}

func (d *Datum) mutexHashFromNamespaces(namespaces []string) int {
	key := d.mutexKeyFromNamespaces(namespaces)

	return d.mutexHashFromKey(key)
}

func (d *Datum) mutexFromHash(hash int) *sync.RWMutex {
	return d.mutexes[hash]
}

func (d *Datum) mutexFromNamespaces(namespaces []string) *sync.RWMutex {
	hash := d.mutexHashFromNamespaces(namespaces)

	return d.mutexFromHash(hash)
}
