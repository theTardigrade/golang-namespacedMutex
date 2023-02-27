package namespacedMutex

import (
	"math"
	"math/big"

	prime "github.com/theTardigrade/golang-prime"
)

const (
	optionsDefaultBucketCount           = 1 << 10
	optionsDefaultMaxUniqueAttemptCount = 1 << 12
)

func (d *Datum) initOptions(opts *Options) {
	if opts == nil {
		opts = new(Options)
	}

	if opts.BucketCount > 0 {
		d.bucketCount = opts.BucketCount
	} else {
		d.bucketCount = optionsDefaultBucketCount
	}

	if opts.BucketCountShouldBePrime && !prime.Is(int64(d.bucketCount)) {
		var exists bool
		var bucketCount64 int64

		bucketCount64, exists = prime.Next(int64(d.bucketCount))
		if !exists || bucketCount64 > math.MaxInt {
			bucketCount64, exists = prime.Prev(int64(d.bucketCount))
		}

		if exists && bucketCount64 <= math.MaxInt {
			d.bucketCount = int(bucketCount64)
		}
	}

	d.bucketCountBig = big.NewInt(int64(d.bucketCount))

	if opts.MaxUniqueAttemptCount > 0 {
		d.maxUniqueAttemptCount = opts.MaxUniqueAttemptCount
	} else {
		d.maxUniqueAttemptCount = optionsDefaultMaxUniqueAttemptCount
	}

	if d.maxUniqueAttemptCount > d.bucketCount {
		d.maxUniqueAttemptCount = d.bucketCount
	}
}
