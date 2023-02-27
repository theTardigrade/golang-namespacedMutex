package namespacedMutex

import (
	"strings"
	"sync"

	hash "github.com/theTardigrade/golang-hash"
)

func (d *Datum) mutexKey(namespaces []string) string {
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

func (d *Datum) mutexHash(key string) int {
	keyHash := hash.Uint64String(key)
	count := uint64(d.mutexesBucketCount)

	return int(keyHash % count)
}

func (d *Datum) mutex(namespaces []string) *sync.RWMutex {
	key := d.mutexKey(namespaces)
	hash := d.mutexHash(key)

	return d.mutexes[hash]
}
