.PHONY: all fmt vet test

all: fmt vet test
	go build -o huang

vet:
	go vet ./...

fmt:
	go fmt ./...

test:
	go test ./... -cover -race
