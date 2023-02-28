package namespacedMutex

// UseRandom allows code to be run within the handler function
// while a mutex based on a randomly generated namespace is
// automatically locked and unlocked before and after use,
// abstracting away the problem of mutual exclusion.
func (d *Datum) UseRandom(isReadOnly bool, handler func()) (err error) {
	mutex, err := d.GetRandomLocked(isReadOnly)
	if err != nil {
		return
	}
	defer mutex.unlock()

	handler()

	return
}
