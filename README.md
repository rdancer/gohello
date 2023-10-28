# gohello

A Go program that reads a message from a file and prints each letter of the message in order, demonstrating the use of goroutines, channels, and mutex locks for concurrent programming in Go.

## How it works

The program consists of three main functions:

1. `main`: Initializes the necessary variables and starts the producer and printer goroutines.

2. `producer`: Reads each letter of the message from the file, along with its index, and sends it to the channel.

3. `printer`: Receives the letters from the channel and prints them in order of their indices.

The program uses a buffered channel to hold the letters and their indices, a mutex lock to synchronize access to shared variables, and a condition variable to wait for the next letter to be printed.

## How to run

1. Clone the repository:

```bash
git clone https://github.com/rdancer/gohello.git
```

2. Navigate to the project directory:

```
cd gohello
```

3. Build, test, and run the program

```
make build
make test
./hello
```
