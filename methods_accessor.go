package namespacedMutex

// GetLocked returns a locked mutex based on the primary and secondary namespaces;
// it must be unlocked after use. The lock will be either read-only or read-write.
func (d *Datum) GetLocked(
	isReadOnly bool,
	primaryNamespace string,
	secondaryNamespaces ...string,
) (mutex *MutexWrapper) {
	rawMutex := d.mutex(primaryNamespace, secondaryNamespaces...)

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
func (d *Datum) Use(handler func(), isReadOnly bool, primaryNamespace string, secondaryNamespaces ...string) {
	mutex := d.GetLocked(isReadOnly, primaryNamespace, secondaryNamespaces...)
	defer mutex.unlock()

	handler()
}
