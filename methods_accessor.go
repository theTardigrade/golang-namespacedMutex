package namespacedMutex

// GetLocked returns a locked mutex based on the given namespaces.
// The mutex must be unlocked after use, and its lock will be
// either read-only or read-write.
func (d *Datum) GetLocked(
	isReadOnly bool,
	namespaces ...string,
) (mutex *MutexWrapper) {
	rawMutex := d.mutexFromNamespaces(namespaces)

	mutex = &MutexWrapper{
		rawMutex:   rawMutex,
		isReadOnly: isReadOnly,
	}

	mutex.lock()

	return
}

const (
	getLockedIfUniqueMaxAttemptCount = 1 << 12
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
	hash := d.mutexHashFromNamespaces(namespaces)
	found = true

	for _, n := range excludedNamespaces {
		if d.mutexHashFromNamespaces(n) == hash {
			found = false
			break
		}
	}

	if !found && d.mutexesBucketCount > 1 {
		attemptCount := getLockedIfUniqueMaxAttemptCount
		if attemptCount > d.mutexesBucketCount {
			attemptCount = d.mutexesBucketCount
		}

		if attemptCount >= 2 {
			excludedHashes := make([]int, len(excludedNamespaces))

			for i, n := range excludedNamespaces {
				excludedHashes[i] = d.mutexHashFromNamespaces(n)
			}

			for i := 2; i <= attemptCount; i++ {
				hash++
				found = true

				for _, h := range excludedHashes {
					if h == hash {
						found = false
						break
					}
				}

				if found {
					break
				}
			}
		}
	}

	if found {
		mutex = &MutexWrapper{
			rawMutex:   d.mutexFromHash(hash),
			isReadOnly: isReadOnly,
		}

		mutex.lock()
	}

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
