.PHONY: build lint tidy check-tidy test cover coverui precommit ci-checks generate install-hooks

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

cover:
	python scripts/check_coverage.py func

coverui:
	python scripts/check_coverage.py html

precommit: build lint test

ci-checks: build lint cover check-tidy

generate:
	go generate ./...
	mockery

install-hooks:
	lefthook install
