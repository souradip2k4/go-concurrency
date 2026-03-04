package main

import (
	"fmt"
	"runtime"
	"sync/atomic"
)

// Asynchronous Preemption Go >= 1.14
func main() {
	// Set the number of logical processors to 1 so parallel execution doesn't happen
	runtime.GOMAXPROCS(1)

	var i int64 = 0

	go func() {
		// atomic.LoadInt64 forces Go to read 'i' from actual memory and not a stale CPU register cache.
		// without atomic it will print 0
		fmt.Println("goroutine ran after i =", atomic.LoadInt64(&i), "th loop iterations")
	}()

	// If we don't use atomic.LoadInt64, it will print 0 because the goroutine will be scheduled before the main goroutine increments i
	// go fmt.Println("goroutine ran after i =", atomic.LoadInt64(&i), "th loop iterations")

	for {
		// atomic.AddInt64 forces every increment to be written to memory and not just kept in a CPU register.
		// Can also use mutex
		atomic.AddInt64(&i, 1)
	}
}
