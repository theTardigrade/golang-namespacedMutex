package namespacedMutex

import (
	"sync"
	"time"

	"github.com/theTardigrade/cache"
)

type Datum struct {
	cache                    *cache.Cache
	masterMutexes            []*sync.Mutex
	masterMutexesBucketCount int
	namespaceSeparator       string
}

type Options struct {
	CacheExpiryDuration      time.Duration
	CacheMaxValues           int
	MasterMutexesBucketCount int
	NamespaceSeparator       string
}

func New(opts *Options) *Datum {
	d := Datum{}

	d.initOptions(opts)

	{
		bc := d.masterMutexesBucketCount

		d.masterMutexes = make([]*sync.Mutex, bc)

		for bc--; bc >= 0; bc-- {
			m := sync.Mutex{}
			d.masterMutexes[bc] = &m
		}
	}

	d.cache = cache.NewCacheWithOptions(cache.Options{
		ExpiryDuration: opts.CacheExpiryDuration,
		MaxValues:      opts.CacheMaxValues,
		PreDeletionFunc: func(key string, value interface{}, setTime time.Time) {
			d.masterMutex(key).Lock()

			if mutex, ok := value.(*sync.RWMutex); ok {
				mutex.Lock()
			}
		},
		PostDeletionFunc: func(key string, value interface{}, setTime time.Time) {
			d.masterMutex(key).Unlock()
		},
	})

	return &d
}
