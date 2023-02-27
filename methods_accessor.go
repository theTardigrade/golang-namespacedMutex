package namespacedMutex

import (
	"strconv"
	"sync"
)

// GetLocked returns a locked mutex based on the given namespaces.
// The mutex must be unlocked after use, and its lock will be
// either read-only or read-write.
func (d *Datum) GetLocked(
	isReadOnly bool,
	namespaces ...string,
) (mutex *MutexWrapper) {
	rawMutex := d.mutex(namespaces)

	mutex = &MutexWrapper{
		rawMutex:   rawMutex,
		isReadOnly: isReadOnly,
	}

	mutex.lock()

	return
}

const (
	getLockedIfUniqueMaxAttempts = 1 << 12
)

// GetLockedIfUnique attempts to return a locked mutex based on the given namespaces.
// However, if any collection of namespaces from the list of excluded namespaces
// produces the same mutex, then no mutex will be returned or locked.
// If a mutex is found, it must be unlocked after use, and its lock will be
// either read-only or read-write.
func (d *Datum) GetLockedIfUnique(
	isReadOnly bool,
	namespaces []string,
	excludedNamespaces [][]string,
) (mutex *MutexWrapper, found bool) {
	rawMutex := d.mutex(namespaces)
	found = true

	for _, n := range excludedNamespaces {
		if d.mutex(n) == rawMutex {
			found = false
			break
		}
	}

	if !found && d.mutexesBucketCount > 1 {
		seenMutexes := make(map[*sync.RWMutex]struct{})

		seenMutexes[rawMutex] = struct{}{}

		for i := 2; i <= getLockedIfUniqueMaxAttempts; i++ {
			namespacesForAttempt := append([]string{
				strconv.FormatInt(int64(i), 36),
			}, namespaces...)

			rawMutex = d.mutex(namespacesForAttempt)
			if _, seen := seenMutexes[rawMutex]; seen {
				continue
			}
			found = true

			for _, n := range excludedNamespaces {
				if d.mutex(n) == rawMutex {
					found = false
					break
				}
			}

			if found {
				break
			}

			seenMutexes[rawMutex] = struct{}{}
		}
	}

	mutex = &MutexWrapper{
		rawMutex:   rawMutex,
		isReadOnly: isReadOnly,
	}

	mutex.lock()

	return
}

// Use allows code to be run within the handler function
// while the mutex is automatically locked and unlocked
// before and after use. It abstracts away the problem
// of mutual exclusion.
func (d *Datum) Use(handler func(), isReadOnly bool, namespaces ...string) {
	mutex := d.GetLocked(isReadOnly, namespaces...)
	defer mutex.unlock()

	handler()
}
