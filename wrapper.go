package namespacedMutex

import "sync"

// MutexWrapper is returned from the GetLocked function;
// it is a wrapper around a read-write mutex.
type MutexWrapper struct {
	rawMutex    *sync.RWMutex
	unlockMutex sync.Mutex
	isReadOnly  bool
	isUnlocked  bool
}

func (m *MutexWrapper) lock() {
	if m.isReadOnly {
		m.rawMutex.RLock()
	} else {
		m.rawMutex.Lock()
	}
}

func (m *MutexWrapper) unlock() {
	if m.isReadOnly {
		m.rawMutex.RUnlock()
	} else {
		m.rawMutex.Unlock()
	}
}

// Unlock must be called when the mutex is no longer in use.
// It can be called multiple times without triggering a panic,
// but it should ideally only be called once after every use.
func (m *MutexWrapper) Unlock() {
	defer m.unlockMutex.Unlock()
	m.unlockMutex.Lock()

	if m.isUnlocked {
		return
	}

	m.unlock()
	m.isUnlocked = true
}

// Raw returns the underlying RWMutex pointer.
// There should ordinarily be no need to call this function.
func (m *MutexWrapper) Raw() *sync.RWMutex {
	return m.rawMutex
}
