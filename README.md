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

var mutexManager = namespacedMutex.New(namespacedMutex.Options{
	MutexesBucketCount:            1 << 5,
	MutexesBucketCountMustBePrime: true,
})

func main() {
	numbers := make([]string, 0, 100)
	var wg sync.WaitGroup

	wg.Add(cap(numbers))

	for i := 1; i <= cap(numbers); i++ {
		go func(i int) {
			defer wg.Done()

			// when the Use function is called, a mutex stored
			// under the namespace will automatically be locked
			// before the handler function runs, and unlocked
			// once it's finished
			mutexManager.Use(false, "this-is-the-namespace", func() {
				numbers = append(numbers, strconv.Itoa(i))
			})
		}(i)
	}

	wg.Wait()
	wg.Add(len(numbers))

	var numbersList string

	for i := 0; i < len(numbers); i++ {
		go func(i int) {
			defer wg.Done()

			// you can also use a mutex directly by calling
			// the GetLocked function
			mutex := mutexManager.GetLocked(false, "another-namespace")
			defer mutex.Unlock()

			numbersList += "(" + numbers[i] + ")"
		}(i)
	}

	wg.Wait()

	fmt.Println(numbersList)
	fmt.Println(len(numbers))
}
```

## Support

If you use this package, or find any value in it, please consider donating:

[![ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/S6S2EIRL0)