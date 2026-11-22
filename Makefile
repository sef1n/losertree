.PHONY: check check-ci fmt fmt-check vet tidy test test-race test-ptr test-cover test-bench test-fuzz-cycle

# Bench params
BTime ?= 1
BCount ?= 1
BOut ?=
BName ?=

# Fuzz params
FTime ?= 1m

check: fmt-check vet test test-race test-ptr

check-ci: check test-fuzz-cycle

fmt:
	gofmt -w .

fmt-check:
	@diff=$$(gofmt -l .); \
	if [ -n "$$diff" ]; then \
		echo "The files are not formatted"; \
		echo "$$diff"; \
		exit 1; \
	fi

vet:
	go vet ./...

tidy:
	go mod tidy -diff

test:
	go test ./...

test-race:
	go test -race ./...

test-ptr:
	go test -gcflags=all=-d=checkptr=2 ./...

test-cover:
	go test -coverprofile=c.out ./...
	go tool cover -func=c.out

test-bench:
	@if [ -z "$(BName)" ]; then \
		echo "Usage: make test-bench BName=BenchmarkSmth"; \
		exit 1; \
	fi
	go test -bench=$(BName) -benchmem -benchtime=$(BTime)x -count=$(BCount) ./... | tee "$(BOut)"

test-fuzz-cycle:
	go test -fuzz=FuzzCycle -fuzztime=$(FTime) .

