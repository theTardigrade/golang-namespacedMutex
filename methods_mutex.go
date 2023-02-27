package namespacedMutex

import (
	"strings"
	"sync"

	hash "github.com/theTardigrade/golang-hash"
)

func (d *Datum) mutexKey(primaryNamespace string, secondaryNamespaces []string) string {
	if len(secondaryNamespaces) == 0 {
		return primaryNamespace
	}

	var builder strings.Builder

	builder.WriteString(primaryNamespace)

	separator := d.namespaceSeparator

	for _, n := range secondaryNamespaces {
		builder.WriteString(separator)
		builder.WriteString(n)
	}

	return builder.String()
}

func (d *Datum) mutexHash(namespace string) int {
	namespaceHash := hash.Uint64String(namespace)
	count := uint64(d.mutexesBucketCount)

	return int(namespaceHash % count)
}

func (d *Datum) mutex(primaryNamespace string, secondaryNamespaces ...string) *sync.RWMutex {
	key := d.mutexKey(primaryNamespace, secondaryNamespaces)
	hash := d.mutexHash(key)

	return d.mutexes[hash]
}
