.PHONY: all
all: fmt test

.PHONY: fmt
fmt:
	go fmt ./...
	go mod tidy

.PHONY: test 
test:
	go test -race ./...
