package namespacedMutex

import (
	"sync"

	cache "github.com/theTardigrade/golang-cache"
)

// Datum is used as the main return type, producing
// namespaced mutexes on demand.
type Datum struct {
	cache                    *cache.Cache
	masterMutexes            []*sync.Mutex
	masterMutexesBucketCount int
	namespaceSeparator       string
}

// Options is used in the New constructor function.
type Options struct {
	MasterMutexesBucketCount            int
	MasterMutexesBucketCountMustBePrime bool
	NamespaceSeparator                  string
}

// New creates a new Datum based on the given options;
// default options will be used, if necessary.
func New(opts Options) *Datum {
	d := Datum{}

	d.initOptions(&opts)

	{
		bc := d.masterMutexesBucketCount

		d.masterMutexes = make([]*sync.Mutex, bc)

		for bc--; bc >= 0; bc-- {
			m := sync.Mutex{}
			d.masterMutexes[bc] = &m
		}
	}

	d.cache = cache.NewCacheWithOptions(cache.Options{
		ExpiryDuration: -1,
		MaxValues:      -1,
	})

	return &d
}

// NewDefault is equivalent to the New constructor
// function with default options.
func NewDefault() *Datum {
	return New(Options{})
}
