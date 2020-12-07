package namespacedMutex

import (
	"strings"
	"sync"
)

func (d *Datum) cacheKey(primaryNamespace string, secondaryNamespaces []string) string {
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

// GetLocked returns a locked mutex based on the primary and secondary namespaces;
// it must be unlocked after use. The lock will be either read-only or read-and-write.
func (d *Datum) GetLocked(readOnly bool, primaryNamespace string, secondaryNamespaces ...string) (mutex *sync.RWMutex) {
	cacheKey := d.cacheKey(primaryNamespace, secondaryNamespaces)
	masterMutex := d.masterMutex(cacheKey)

	defer masterMutex.Unlock()
	masterMutex.Lock()

	defer func() {
		if readOnly {
			mutex.RLock()
		} else {
			mutex.Lock()
		}
	}()

	if mutexInterface, ok := d.cache.Get(cacheKey); ok {
		if mutex, ok = mutexInterface.(*sync.RWMutex); ok {
			return
		}
	}

	mutex = &sync.RWMutex{}

	d.cache.Set(cacheKey, mutex)

	return
}
