package namespacedMutex

import (
	"sync"
	"time"

	"github.com/theTardigrade/cache"
)

type Datum struct {
	masterMutexesBucketCount int
	masterMutexes            []*sync.Mutex
	cache                    *cache.Cache
	namespaceSeparator       string
}

type Options struct {
	MasterMutexesBucketCount int
	CacheExpiryDuration      time.Duration
	CacheMaxValues           int
	NamespaceSeparator       string
}

func New(opts *Options) *Datum {
	d := Datum{}

	d.initOptions(opts)

	d.cache = cache.NewCache(opts.CacheExpiryDuration, opts.CacheMaxValues)

	bc := d.masterMutexesBucketCount

	d.masterMutexes = make([]*sync.Mutex, bc)

	for bc--; bc >= 0; bc-- {
		m := sync.Mutex{}
		d.masterMutexes[bc] = &m
	}

	return &d
}
