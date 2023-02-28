package namespacedMutex

// Use allows code to be run within the handler function
// while a mutex based on the given namespace is
// automatically locked and unlocked before and after use,
// abstracting away the problem of mutual exclusion.
func (d *Datum) Use(isReadOnly bool, namespace string, handler func()) {
	mutex := d.GetLocked(isReadOnly, namespace)
	defer mutex.unlock()

	handler()
}
