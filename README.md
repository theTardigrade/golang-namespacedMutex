# golang-namespacedMutex

## Support

If you use or appreciate this package in any way, please consider donating at [PayPal](https://www.paypal.me/jismithpp).

## Example

```golang
package main

import (
	"fmt"

	namespacedMutex "github.com/theTardigrade/golang-namespacedMutex"
)

var (
	m = namespacedMutex.New(namespacedMutex.Options{})
	n int
)

const iterations = 2e4

func main() {
	c := make(chan struct{})

	for i := 1; i <= iterations; i++ {
		go func(i int) {
			// a mutex stored under the given namespace will automatically be
			// locked before and unlocked after the given handler function runs
			m.Use(func() {
				if n%2 == 0 {
					n++
				} else {
					n += i
				}
			}, false, "main")

			// notify the main function that this goroutine has completed its work
			c <- struct{}{}
		}(i)
	}

	for i := 0; i < iterations; i++ {
		<-c
	}

	fmt.Printf("MAGIC NUMBER: %d\n", n)
}
```