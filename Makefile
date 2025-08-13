BINARY=gosort
BIN_DIR=bin

.PHONY: build test lint run clean

build:
	@mkdir -p $(BIN_DIR)
	go build -o ${BIN_DIR}/${BINARY} ./cmd/gosort

run: build
	./${BIN_DIR}/${BINARY}

test:
	go test -v ./...
	go test -v -tags=integration

lint:
	go vet ./...
	golangci-lint run ./...

clean:
	rm -rf ${BIN_DIR}
	rm -f gosort