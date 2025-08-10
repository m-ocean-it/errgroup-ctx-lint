.PHONY: build run test clean

SRC := $(shell find . -type f -name '*.go' ! -name '*_test.go')

build: ./bin/errgroup-ctx-lint

run: build
	./bin/errgroup-ctx-lint

test:
	go test ./...

clean:
	rm -rf ./bin

./bin/errgroup-ctx-lint: $(SRC) | ./bin
	go build -o ./bin/errgroup-ctx-lint ./cmd/ErrGroupCtxLint

./bin:
	mkdir -p ./bin
