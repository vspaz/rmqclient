BINARY_NAME=rmq

all: build
build:
	go build -o $(BINARY_NAME) examples/producer/main.go

.PHONY: clean
clean:
	rm -f $(BINARY_NAME)

.PHONY: test
test:
	go test -race -v -p 4 ./...

.PHONY: style-fix
style-fix:
	gofmt -w .

.PHONY: lint
lint:
	golangci-lint run
