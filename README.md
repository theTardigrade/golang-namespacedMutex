# golang-namespacedMutex

## Example

```golang
package main

import (
	"fmt"
	"strconv"
	"sync"

	namespacedMutex "github.com/theTardigrade/golang-namespacedMutex"
)

var (
	mutexManager = namespacedMutex.New(namespacedMutex.Options{
		CacheExpiryDuration:                 -1,
		CacheMaxValues:                      -1,
		MasterMutexesBucketCount:            1 << 20,
		MasterMutexesBucketCountMustBePrime: true,
	})
)

const iterations = 21

func main() {
	var numbers []string
	var wg sync.WaitGroup

	wg.Add(iterations)

	for i := 1; i <= iterations; i++ {
		go func(i int) {
			defer wg.Done()

			// when the Use function is called, a mutex stored
			// under the namespace will automatically be locked
			// before the handler function runs, and unlocked
			// once it's finished
			mutexManager.Use(func() {
				numbers = append(numbers, strconv.Itoa(i))
			}, false, "this-is-the-namespace")
		}(i)
	}

	wg.Wait()

	fmt.Println(numbers)
	fmt.Println(len(numbers))
}
```