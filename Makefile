.PHONY: check check-ci fmt fmt-check vet tidy test test-race test-ptr test-cover test-bench test-fuzz-cycle

SHELL := /bin/bash
.SHELLFLAGS := -euo pipefail -c

# Bench params
BTime ?= 1x
BCount ?= 1
BOut ?= /dev/null
BName ?=
BPkgs ?= ./benchmarks/losertree/ ./benchmarks/goloser/

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
		echo "Usage: make test-bench BName=BenchmarkSmth BPkgs=./benchmarks/losertree/"; \
		exit 1; \
	fi
	BENCHNAME="$(BName)" BENCHTIME=$(BTime) COUNT=$(BCount) ./benchmarks/gobench.sh $(BPkgs) | tee "$(BOut)"

test-fuzz-cycle:
	go test -fuzz=FuzzCycle -fuzztime=$(FTime) .

