package namespacedMutex

import "time"

const (
	optionsDefaultBucketCount         = 1 << 10
	optionsDefaultCacheMaxValues      = 1 << 16
	optionsDefaultCacheExpiryDuration = time.Hour
	optionsDefaultNamespaceSeparator  = "__::__"
)

const (
	optionsMaxBucketCount = 1 << 30
)

func (d *Datum) initOptions(opts *Options) {
	if opts == nil {
		*opts = Options{}
	}

	if opts.MasterMutexesBucketCount > 0 {
		if opts.MasterMutexesBucketCount <= optionsMaxBucketCount {
			d.masterMutexesBucketCount = opts.MasterMutexesBucketCount
		} else {
			d.masterMutexesBucketCount = optionsMaxBucketCount
		}
	} else {
		d.masterMutexesBucketCount = optionsDefaultBucketCount
	}

	if opts.CacheMaxValues == 0 {
		opts.CacheMaxValues = optionsDefaultCacheMaxValues
	}

	if opts.CacheExpiryDuration == 0 {
		opts.CacheExpiryDuration = optionsDefaultCacheExpiryDuration
	}

	if opts.NamespaceSeparator == "" {
		d.namespaceSeparator = optionsDefaultNamespaceSeparator
	} else {
		d.namespaceSeparator = opts.NamespaceSeparator
	}
}
