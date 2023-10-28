package main

import (
	"fmt"
	"sync"
	"time"
)

const message = "Hello, world!\n"
const RunThisManyInParallel = 400
const PollInterval = time.Millisecond * 5

func main() {
	var wg sync.WaitGroup
	ch := make(chan struct {
		letter rune
		index  int
	})
	nextIndex := 0
	mu := sync.Mutex{}
	cond := sync.NewCond(&mu)

	// Start producer goroutine
	go producer(ch)

	// Increment the WaitGroup counter
	wg.Add(RunThisManyInParallel)

	// Start printer goroutines
	for i := 0; i < RunThisManyInParallel; i++ {
		go printer(&wg, ch, &mu, &nextIndex, cond)
	}

	// Wait for all printer goroutines to finish
	wg.Wait()
}

// producer sends each letter of the message, along with its index, to the channel
func producer(ch chan<- struct {
	letter rune
	index  int
}) {
	for index, letter := range message {
		ch <- struct {
			letter rune
			index  int
		}{letter, index}
	}
	close(ch)
}

// printer prints the letters from the channel in the order of their indices
func printer(wg *sync.WaitGroup, ch <-chan struct {
    letter rune
    index  int
}, mu *sync.Mutex, nextIndex *int, cond *sync.Cond) {
    defer wg.Done()
    for {
        letterWithIndex, ok := <-ch
        if !ok {
            return
        }
	mu.Lock() // acquire lock
        // Wait until it's the turn of this letter to be printed
        for letterWithIndex.index != *nextIndex {
	    fmt.Print(".")
            cond.Wait() // release lock temporarily and suspend thread
        }
	// Print the letter and mark its index as printed
        fmt.Print(string(letterWithIndex.letter))
        *nextIndex++
	cond.Broadcast() // wake up other threads; they will all race for the lock, one of them will acquire it, the other ones will continue to Wait()
        mu.Unlock() // release the lock
    }
}

