package namespacedMutex

import (
	prime "github.com/theTardigrade/golang-prime"
)

const (
	optionsDefaultBucketCount        = 1 << 10
	optionsDefaultNamespaceSeparator = "__::__"
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

	if opts.MasterMutexesBucketCountMustBePrime {
		d.masterMutexesBucketCount = prime.Next(d.masterMutexesBucketCount)
	}

	if opts.NamespaceSeparator == "" {
		d.namespaceSeparator = optionsDefaultNamespaceSeparator
	} else {
		d.namespaceSeparator = opts.NamespaceSeparator
	}
}
