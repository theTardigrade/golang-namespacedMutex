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

// GetLockedMutex is returned from the GetLocked function;
// it is a wrapper around a read-write mutex.
type GetLockedMutex struct {
	realMutex   *sync.RWMutex
	unlockMutex sync.Mutex
	isReadOnly  bool
	isUnlocked  bool
}

func (g *GetLockedMutex) lock() {
	if g.isReadOnly {
		g.realMutex.RLock()
	} else {
		g.realMutex.Lock()
	}
}

func (g *GetLockedMutex) unlock() {
	if g.isReadOnly {
		g.realMutex.RUnlock()
	} else {
		g.realMutex.Unlock()
	}
}

// Unlock must be called when the mutex is no longer in use.
// It can be called multiple times without triggering a panic,
// but it should ideally only be called once after every use.
func (g *GetLockedMutex) Unlock() {
	defer g.unlockMutex.Unlock()
	g.unlockMutex.Lock()

	if g.isUnlocked {
		return
	}

	g.unlock()
	g.isUnlocked = true
}

// Raw returns the underlying RWMutex pointer.
// There should ordinarily be no need to call this function.
func (g *GetLockedMutex) Raw() *sync.RWMutex {
	return g.realMutex
}

// GetLocked returns a locked mutex based on the primary and secondary namespaces;
// it must be unlocked after use. The lock will be either read-only or read-write.
func (d *Datum) GetLocked(isReadOnly bool, primaryNamespace string, secondaryNamespaces ...string) (mutex *GetLockedMutex) {
	cacheKey := d.cacheKey(primaryNamespace, secondaryNamespaces)
	masterMutex := d.masterMutex(cacheKey)

	defer masterMutex.Unlock()
	masterMutex.Lock()

	var realMutex *sync.RWMutex

	defer func() {
		mutex = &GetLockedMutex{
			realMutex:  realMutex,
			isReadOnly: isReadOnly,
		}
		mutex.lock()
	}()

	if realMutexInterface, ok := d.cache.Get(cacheKey); ok {
		if realMutex, ok = realMutexInterface.(*sync.RWMutex); ok {
			return
		}
	}

	realMutex = &sync.RWMutex{}

	d.cache.Set(cacheKey, realMutex)

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
