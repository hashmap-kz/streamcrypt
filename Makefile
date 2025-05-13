.PHONY: lint
lint:
	golangci-lint run --output.tab.path=stdout

.PHONY: test
test:
	go test -v -cover ./...

.PHONY: test-cov
test-cov:
	go test -coverprofile=coverage.txt ./...
	go tool cover -html=coverage.txt
