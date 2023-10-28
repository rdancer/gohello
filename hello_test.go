package main

import (
	"bytes"
	"os"
	"sync"
	"testing"
)

var MFP string // MessageFilePath

func TestProducer(t *testing.T) {
	ch := make(chan struct {
		letter rune
		index  int
	})
	go producer(ch, MFP)

	var buf bytes.Buffer
	for i := 0; i < len(message); i++ {
		item := <-ch
		if item.letter != rune(message[i]) || item.index != i {
			t.Errorf("producer sent incorrect item: got (%v, %v), want (%v, %v)", item.letter, item.index, message[i], i)
		}
		buf.WriteRune(item.letter)
	}
	if buf.String() != message {
		t.Errorf("producer sent incorrect message: got %q, want %q", buf.String(), message)
	}
}

func TestPrinter(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	ch := make(chan struct {
		letter rune
		index  int
	})
	go func() {
		for i := 0; i < len(message); i++ {
			ch <- struct {
				letter rune
				index  int
			}{rune(message[i]), i}
		}
		close(ch)
	}()

	// open a temporary file
	tmpfile, err := os.CreateTemp("", "*.output.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up
	// redirect stdout to the temporary file
	stdout := os.Stdout
	os.Stdout = tmpfile
	defer func() { os.Stdout = stdout }() // clean up

	mu := &sync.Mutex{}
	nextIndex := 0
	cond := sync.NewCond(mu)
	go printer(&wg, ch, mu, &nextIndex, cond)

	wg.Wait()
	// read the temporary file into a string
	if _, err := tmpfile.Seek(0, 0); err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(tmpfile); err != nil {
		t.Fatal(err)
	}
	if buf.String() != message {
		t.Errorf("printer printed incorrect message: got %q, want %q", buf.String(), message)
	}
}

func TestMain(m *testing.M) {
	// create a temporary message file for testing
	tmpfile, err := os.CreateTemp("", "*.message.txt")
	if err != nil {
		panic(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up
	if _, err := tmpfile.Write([]byte(message)); err != nil {
		panic(err)
	}
	if err := tmpfile.Close(); err != nil {
		panic(err)
	}
	MFP = tmpfile.Name()

	// run tests
	exitCode := m.Run()

	os.Exit(exitCode)
}

const message = "Test Test Test!\n"
