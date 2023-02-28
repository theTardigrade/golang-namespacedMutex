package namespacedMutex

import (
	"crypto/rand"
)

// GetRandomLocked returns a locked mutex based on a randomly
// generated namespace.
// The mutex must be unlocked after use, and its lock will be
// either read-only or read-write.
func (d *Datum) GetRandomLocked(isReadOnly bool) (mutex *MutexWrapper, err error) {
	indexBig, err := rand.Int(rand.Reader, d.bucketCountBig)
	if err != nil {
		return
	}

	rawMutex := d.mutexFromIndex(int(indexBig.Uint64()))
	mutex = newMutexWrapper(rawMutex, isReadOnly, true)

	return
}

// GetRandomLockedIfUnique attempts to return a locked mutex based on
// a randomly generated namespace.
// However, if any of the comparison namespaces give the same mutex,
// then no mutex will be returned or locked.
// A number of other primary keys, related to the first, will be attempted
// before the search for a unique mutex is ended.
// If a mutex is found, it must be unlocked after use, and its lock will be
// either read-only or read-write.
func (d *Datum) GetRandomLockedIfUnique(
	isReadOnly bool,
	comparisonNamespaces ...string,
) (mutex *MutexWrapper, found bool, err error) {
	randomNamespaceBytes := make([]byte, 1024)

	_, err = rand.Read(randomNamespaceBytes)
	if err != nil {
		return
	}

	mutex, found = d.GetLockedIfUnique(
		isReadOnly,
		string(randomNamespaceBytes),
		comparisonNamespaces...,
	)

	return
}
