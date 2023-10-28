build:
	go build -o hello

run:
	make build && ./hello

test:
	go test ./... -v
