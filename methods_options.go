package namespacedMutex

import (
	"math"
	"math/big"

	prime "github.com/theTardigrade/golang-prime"
)

const (
	optionsDefaultBucketCount = 1 << 10
)

const (
	optionsMaxBucketCount = math.MaxInt
)

func (d *Datum) initOptions(opts *Options) {
	if opts == nil {
		*opts = Options{}
	}

	if opts.MutexesBucketCount > 0 {
		if opts.MutexesBucketCount <= optionsMaxBucketCount {
			d.mutexesBucketCount = opts.MutexesBucketCount
		} else {
			d.mutexesBucketCount = optionsMaxBucketCount
		}
	} else {
		d.mutexesBucketCount = optionsDefaultBucketCount
	}

	if opts.MutexesBucketCountMustBePrime && !prime.Is(int64(d.mutexesBucketCount)) {
		var exists bool
		var bucketCount64 int64

		bucketCount64, exists = prime.Next(int64(d.mutexesBucketCount))
		if !exists || bucketCount64 > math.MaxInt {
			bucketCount64, exists = prime.Prev(int64(d.mutexesBucketCount))
		}

		if exists && bucketCount64 <= math.MaxInt {
			d.mutexesBucketCount = int(bucketCount64)
		}
	}

	d.mutexesBucketCountBig = big.NewInt(int64(d.mutexesBucketCount))
}
