# golang-namespacedMutex

This Go package enables a theoretically infinite number of mutexes to be stored and accessed based on namespaces.

[![Go Reference](https://pkg.go.dev/badge/github.com/theTardigrade/golang-namespacedMutex/.svg)](https://pkg.go.dev/github.com/theTardigrade/golang-namespacedMutex/) [![Go Report Card](https://goreportcard.com/badge/github.com/theTardigrade/golang-namespacedMutex)](https://goreportcard.com/report/github.com/theTardigrade/golang-namespacedMutex)

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
		MasterMutexesBucketCount:            1 << 5,
		MasterMutexesBucketCountMustBePrime: true,
	})
)

const iterations = 100

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

## Support

If you use this package, or find any value in it, please consider donating:

[![ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/S6S2EIRL0)