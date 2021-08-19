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
// it must be unlocked after use. The lock will be either read-only or read-write.
func (d *Datum) GetLocked(isReadOnly bool, primaryNamespace string, secondaryNamespaces ...string) (mutex *MutexWrapper) {
	cacheKey := d.cacheKey(primaryNamespace, secondaryNamespaces)
	masterMutex := d.masterMutex(cacheKey)

	defer masterMutex.Unlock()
	masterMutex.Lock()

	var rawMutex *sync.RWMutex

	defer func() {
		mutex = &MutexWrapper{
			rawMutex:   rawMutex,
			isReadOnly: isReadOnly,
		}
		mutex.lock()
	}()

	if realMutexInterface, ok := d.cache.Get(cacheKey); ok {
		if rawMutex, ok = realMutexInterface.(*sync.RWMutex); ok {
			return
		}
	}

	rawMutex = &sync.RWMutex{}

	d.cache.Set(cacheKey, rawMutex)

	return
}

// Use allows code to be run within the handler function
// while the mutex is automatically locked and unlocked
// before and after use. It abstracts away the problem
// of mutual exclusion.
func (d *Datum) Use(handler func(), isReadOnly bool, primaryNamespace string, secondaryNamespaces ...string) {
	mutex := d.GetLocked(isReadOnly, primaryNamespace, secondaryNamespaces...)
	defer mutex.unlock()

	handler()
}
