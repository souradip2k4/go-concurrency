package main

import (
	"runtime"
	"sync/atomic"
)

func main() {
	runtime.GOMAXPROCS(1)
	// Possible with go version >= 1.14
	var i int64 = 0
	// Both workers get stuck
	go func() {
		for {
			atomic.AddInt64(&i, 1)
		}
	}() // Occupies Worker 1
	go func() {
		for {
			atomic.AddInt64(&i, 1)
		}
	}() // Occupies Worker 2

	// Now both "lanes" are full!
	// Wrap in a closure so atomic.LoadInt64 is called when the goroutine
	// actually RUNS (after preemption), not when it is scheduled (i=0).
	go func() {
		println("Can I run after i =", atomic.LoadInt64(&i), "th loop iterations")
	}()

	for {
		atomic.AddInt64(&i, 1)
	} // Occupies Worker 1 or 2 depending on where main is
}
