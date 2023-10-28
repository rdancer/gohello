package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

const MessageFilePath = "./message.txt"
const RunThisManyInParallel = 4

func main() {
	var wg sync.WaitGroup
	ch := make(chan struct {
		letter rune
		index  int
	}, 1000) // buffered channel with a capacity of 1000
	readyCh := make(chan struct{}, 1)
	readyCh <- struct{}{}
	nextIndex := 0
	mu := sync.Mutex{}
	cond := sync.NewCond(&mu)

	// Start producer goroutine
	go producer(ch, readyCh, MessageFilePath)

	// Increment the WaitGroup counter
	wg.Add(RunThisManyInParallel)

	// Start printer goroutines
	counter := RunThisManyInParallel
	for i := 0; i < RunThisManyInParallel; i++ {
		go func() {
		    printer(&wg, ch, &mu, &nextIndex, cond, readyCh)
		    fmt.Println("printer done")
		    mu.Lock()
		    counter--
		    mu.Unlock()
		}()
	}

	// Wait for all printer goroutines to finish
	for {
	    mu.Lock() // not really necessary; placate the race detector
	    if counter > 0 {
		mu.Unlock()
		fmt.Println("counter:", counter)
		cond.Broadcast()
		time.Sleep(time.Second)
	    } else {
		mu.Unlock()
	    }
	}
	wg.Wait() // no-op

	close(readyCh) // Close readyCh channel after all goroutines are done
}

// producer sends each letter of the message, along with its index, to the channel
func producer(ch chan<- struct {
    letter rune
    index  int
}, readyCh <-chan struct{}, filename string) {
    file, err := os.Open(filename)
    if err != nil {
        fmt.Fprintln(os.Stderr, "Error opening file:", err)
        return
    }
    defer file.Close()

    index := 0
    reader := bufio.NewReader(file)
    for {
        _, ok := <-readyCh
        if !ok {
            break
        }

        r, _, err := reader.ReadRune()
        if err == io.EOF {
            break
        }
        if err != nil {
            fmt.Fprintln(os.Stderr, "Error reading rune from file:", err)
            return
        }
        ch <- struct {
            letter rune
            index  int
        }{r, index}
        index++
    }
    close(ch)
}


// printer prints the letters from the channel in the order of their indices
func printer(wg *sync.WaitGroup, ch <-chan struct {
    letter rune
    index  int
}, mu *sync.Mutex, nextIndex *int, cond *sync.Cond, readyCh chan<- struct{}) {
    defer wg.Done()
    for {
        // Send signal that we are ready for more data
        readyCh <- struct{}{}

        letterWithIndex, ok := <-ch
        if !ok {
            return
        }

	mu.Lock() // acquire lock
        // Wait until it's the turn of this letter to be printed
        for letterWithIndex.index != *nextIndex {
	    fmt.Print(".")
            cond.Wait() // release lock temporarily and suspend thread
	    fmt.Println("Wakey-wakey")
        }
	// Print the letter and mark its index as printed
        fmt.Print(string(letterWithIndex.letter))
        *nextIndex++
	cond.Broadcast() // wake up other threads; they will all race for the lock, one of them will acquire it, the other ones will continue to Wait()
        mu.Unlock() // release the lock
    }
}

