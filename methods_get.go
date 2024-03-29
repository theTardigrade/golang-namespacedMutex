package namespacedMutex

// GetLocked returns a locked mutex based on the given namespace.
// The mutex must be unlocked after use, and its lock will be
// either read-only or read-write.
func (d *Datum) GetLocked(
	isReadOnly bool,
	namespace string,
) (mutex *MutexWrapper) {
	rawMutex := d.mutexFromNamespace(namespace)
	mutex = newMutexWrapper(rawMutex, isReadOnly, true)

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
	index := d.mutexIndexFromNamespace(namespace)
	found = true

	if len(comparisonNamespaces) > 0 {
		for _, n := range comparisonNamespaces {
			if d.mutexIndexFromNamespace(n) == index {
				found = false
				break
			}
		}

		if !found && d.maxUniqueAttemptCount > 1 {
			comparisonIndexes := make([]int, len(comparisonNamespaces))

			for i, n := range comparisonNamespaces {
				comparisonIndexes[i] = d.mutexIndexFromNamespace(n)
			}

			for i := 2; i <= d.maxUniqueAttemptCount; i++ {
				if index == d.bucketCount-1 {
					index = 0
				} else {
					index++
				}

				found = true

				for _, h := range comparisonIndexes {
					if h == index {
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
		rawMutex := d.mutexFromIndex(index)
		mutex = newMutexWrapper(rawMutex, isReadOnly, true)
	}

	return
}
