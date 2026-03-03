package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// try running this example with the -race flag
// cd into "mutex-implementation"
// go run -race main.go
func main() {
	var count int
	var mu mutex
	var wg sync.WaitGroup

	// if we increase the number of go routines
	// this could quickly be detected as a race condition
	// due to too many go routines' status being awake
	// Also there's no way we can control the go routines queue
	// or have access to the runtime internals
	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		go func() {
			defer wg.Done()
			mu.Lock()
			count++
			mu.Unlock()
		}()
	}
	wg.Wait()

	fmt.Println("count", count)
}

type mutex struct {
	state int32
}

func (mu *mutex) Lock() {
	if atomic.CompareAndSwapInt32(&mu.state, 0, 1) {
		return
	}

	// Slow path: the lock is held, so we enter a loop to wait/check
	for {
		// Attempt to signal our presence by incrementing the state.
		// NOTE: In this specific implementation, state > 1 is treated as a failure.
		atomic.AddInt32(&mu.state, 1)
		s := atomic.LoadInt32(&mu.state)

		if s > 1 {
			// This implementation panics if there is ANY contention (waiters > 0)
			panic("all goroutines are asleep - deadlock!")
		}
		// If the state is 1, it means that the mutex is locked
		if s == 1 {
			return
		}
	}
}

func (mu *mutex) Unlock() {
	for atomic.CompareAndSwapInt32(&mu.state, 1, 0) {
		return
	}
	// If the state is not 1, it means that the mutex is unlocked
	panic("unlock of unlocked mutex")
}

/*
func (mu *mutex) Lock() {
	// Keep trying to change state from 0 (unlocked) to 1 (locked)
	for !atomic.CompareAndSwapInt32(&mu.state, 0, 1) {
		// Yield the processor to let other goroutines run
		// This prevents "deadlocking" the CPU while waiting
		runtime.Gosched()
	}
}

func (mu *mutex) Unlock() {
	// Change state from 1 back to 0
	if !atomic.CompareAndSwapInt32(&mu.state, 1, 0) {
		panic("unlock of unlocked mutex")
	}
}
*/
