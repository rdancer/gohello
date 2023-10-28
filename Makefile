build:
	go build -o hello

run:
	@echo --------------------------------------------------------------
	make build && ./hello

test:
	go test ./... -v
