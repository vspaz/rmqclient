PRODUCER_NAME=producer
CONSUMER_NAME=consumer

build-producer:
	go build -ldflags="-s -w" -o $(PRODUCER_NAME) examples/producer/main.go; upx producer

build-consumer:
	go build -ldflags="-s -w" -o $(CONSUMER_NAME) examples/consumer/main.go; upx consumer

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

.PHONY: upgrade
upgrade:
	go mod tidy
	go get -u all ./...