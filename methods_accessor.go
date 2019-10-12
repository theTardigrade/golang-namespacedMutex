package namespacedMutex

import (
	"strings"
	"sync"
)

func (d *Datum) key(primaryNamespace string, secondaryNamespaces []string) string {
	var builder strings.Builder

	builder.WriteString(primaryNamespace)

	if len(secondaryNamespaces) > 0 {
		separator := d.namespaceSeparator

		for _, n := range secondaryNamespaces {
			builder.WriteString(separator)
			builder.WriteString(n)
		}
	}

	return builder.String()
}

func (d *Datum) Get(primaryNamespace string, secondaryNamespaces ...string) *sync.Mutex {
	key := d.key(primaryNamespace, secondaryNamespaces)
	masterMutex := d.masterMutex(key)

	defer masterMutex.Unlock()
	masterMutex.Lock()

	if mutexInterface, ok := d.cache.Get(key); ok {
		if mutex, ok := mutexInterface.(*sync.Mutex); ok {
			return mutex
		}
	}

	mutex := &sync.Mutex{}

	d.cache.Set(key, mutex)

	return mutex
}
