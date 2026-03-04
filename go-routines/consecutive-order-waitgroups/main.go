package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {
	// Set the number of logical processors to 1 so parallel execution doesn't happen
	// True concurrency
	runtime.GOMAXPROCS(1)
	start := time.Now()
	defer func() {
		fmt.Printf("total time: %v\n", time.Since(start))
	}()

	wg.Add(1)
	go task1()
	// wg.Wait() // Uncommenting would cause sequential execution

	wg.Add(1)
	go task2()
	// wg.Wait()

	wg.Add(1)
	go task3()
	wg.Wait()
}

func task1() {
	time.Sleep(100 * time.Millisecond)
	fmt.Println("task 1")
	wg.Done()
}

func task2() {
	time.Sleep(50 * time.Millisecond)
	fmt.Println("task 2")
	wg.Done()
}

func task3() {
	time.Sleep(10 * time.Millisecond)
	fmt.Println("task 3")
	wg.Done()
}
