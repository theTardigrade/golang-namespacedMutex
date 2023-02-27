package namespacedMutex

import (
	"math/big"
	"sync"
)

// Datum is used as the main return type, producing
// namespaced mutexes on demand.
type Datum struct {
	mutexes               []*sync.RWMutex
	mutexesBucketCount    int
	mutexesBucketCountBig *big.Int
}

// Options is used in the New constructor function.
type Options struct {
	MutexesBucketCount            int
	MutexesBucketCountMustBePrime bool
}

// New creates a new Datum based on the given options;
// default options will be used, if necessary.
func New(opts Options) *Datum {
	d := Datum{}

	d.initOptions(&opts)

	{
		bc := d.mutexesBucketCount

		d.mutexes = make([]*sync.RWMutex, bc)

		for bc--; bc >= 0; bc-- {
			d.mutexes[bc] = new(sync.RWMutex)
		}
	}

	return &d
}

// NewDefault is equivalent to the New constructor
// function with default options.
func NewDefault() *Datum {
	return New(Options{})
}
