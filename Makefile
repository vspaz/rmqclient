TARGET_PRODUCER=producer
TARGET_CONSUMER=consumer

LDFLAGS="-s -w"

build-producer:
	go build -ldflags=$(LDGLAGS) -o $(TARGET_PRODUCER) examples/producer/main.go; upx $(TARGET_PRODUCER)

build-consumer:
	go build -ldflags=$(LDFLAGS) -o $(TARGET_CONSUMER) examples/consumer/main.go; upx $(TARGET_CONSUMER)

.PHONY: clean
clean:
	rm -f $(TARGET_PRODUCER) $(TARGET_CONSUMER)

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
