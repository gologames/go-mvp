.PHONY: build lint tidy check-tidy test precommit

build:
	go build ./...

lint:
	golangci-lint run

tidy:
	go mod tidy
	git diff --quiet go.mod go.sum || git add go.mod go.sum

check-tidy:
	go mod tidy
	git diff --exit-code go.mod go.sum || (echo "Run 'make tidy'"; exit 1)

test:
	go test ./...

precommit: build lint check-tidy test
