package namespacedMutex

import (
	"strings"
	"sync"
)

func (d *Datum) cacheKey(primaryNamespace string, secondaryNamespaces []string) string {
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

func (d *Datum) Get(primaryNamespace string, secondaryNamespaces ...string) (mutex *sync.Mutex) {
	cacheKey := d.cacheKey(primaryNamespace, secondaryNamespaces)
	masterMutex := d.masterMutex(cacheKey)

	defer masterMutex.Unlock()
	masterMutex.Lock()

	if mutexInterface, ok := d.cache.Get(cacheKey); ok {
		if mutex, ok = mutexInterface.(*sync.Mutex); ok {
			return
		}
	}

	mutex = &sync.Mutex{}

	d.cache.Set(cacheKey, mutex)

	return
}
