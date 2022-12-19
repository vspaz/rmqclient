.PHONY: test
test:
	go test -race -v -p 4 ./...

.PHONY: style-fix
style-fix:
	gofmt -w .

.PHONY: lint
lint:
	golangci-lint run
