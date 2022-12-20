PRODUCER_NAME=producer
CONSUMER_NAME=consumer

build-producer:
	go build -o $(PRODUCER_NAME) examples/producer/main.go

build-consumer:
	go build -o $(CONSUMER_NAME) examples/consumer/main.go

.PHONY: clean
clean:
	rm -f $(PRODUCER_NAME) $(CONSUMER_NAME)

.PHONY: test
test:
	go test -race -v -p 4 ./...

.PHONY: style-fix
style-fix:
	gofmt -w .

.PHONY: lint
lint:
	golangci-lint run
