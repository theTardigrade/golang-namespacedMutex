package namespacedMutex

// GetLocked returns a locked mutex based on the given namespace.
// The mutex must be unlocked after use, and its lock will be
// either read-only or read-write.
func (d *Datum) GetLocked(
	isReadOnly bool,
	namespace string,
) (mutex *MutexWrapper) {
	rawMutex := d.mutexFromNamespace(namespace)

	mutex = &MutexWrapper{
		rawMutex:   rawMutex,
		isReadOnly: isReadOnly,
	}

	mutex.lock()

	return
}

// GetLockedIfUnique attempts to return a locked mutex based on the given namespace.
// However, if any of the comparison namespaces give the same mutex,
// then no mutex will be returned or locked.
// A number of other primary keys, related to the first, will be attempted
// before the search for a unique mutex is ended.
// If a mutex is found, it must be unlocked after use, and its lock will be
// either read-only or read-write.
func (d *Datum) GetLockedIfUnique(
	isReadOnly bool,
	namespace string,
	comparisonNamespaces ...string,
) (mutex *MutexWrapper, found bool) {
	hash := d.mutexHashFromNamespace(namespace)
	found = true

	if len(comparisonNamespaces) > 0 {
		for _, n := range comparisonNamespaces {
			if d.mutexHashFromNamespace(n) == hash {
				found = false
				break
			}
		}

		if !found && d.maxUniqueAttemptCount > 1 {
			comparisonHashes := make([]int, len(comparisonNamespaces))

			for i, n := range comparisonNamespaces {
				comparisonHashes[i] = d.mutexHashFromNamespace(n)
			}

			for i := 2; i <= d.maxUniqueAttemptCount; i++ {
				if hash == d.bucketCount-1 {
					hash = 0
				} else {
					hash++
					hash %= d.bucketCount
				}

				found = true

				for _, h := range comparisonHashes {
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
func (d *Datum) Use(isReadOnly bool, namespace string, handler func()) {
	mutex := d.GetLocked(isReadOnly, namespace)
	defer mutex.unlock()

	handler()
}
