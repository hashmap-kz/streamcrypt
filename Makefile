# Variables
COV_REPORT 	:= coverage.txt
TEST_FLAGS 	:= -v -race -timeout 30s

# Lint the code
.PHONY: lint
lint:
	golangci-lint run --output.tab.path=stdout

# Run unit tests
.PHONY: test
test:
	go test -v -cover ./...

# Run tests with coverage
.PHONY: test-cov
test-cov:
	go test -coverprofile=$(COV_REPORT) ./...
	go tool cover -html=$(COV_REPORT)
