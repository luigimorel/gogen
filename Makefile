BINARY_NAME=gogen
BINARY_PATH=./bin/$(BINARY_NAME)
MAIN_PATH=./main.go
TMP_BINARY=./tmp/$(BINARY_NAME)

.PHONY: all build build-dev run dev test test-verbose test-coverage fmt lint vet check deps tidy clean install build-all

all: build

build:
	@mkdir -p bin
	go build -o $(BINARY_PATH) $(MAIN_PATH)

build-dev:
	@mkdir -p tmp
	go build -o $(TMP_BINARY) .

run: build
	$(BINARY_PATH)

dev:
	air

test:
	go test ./...

test-verbose:
	go test -v ./...

test-coverage:
	go test -cover ./...

fmt:
	go fmt ./...

lint:
	golangci-lint run

vet:
	go vet ./...

check: fmt vet lint

deps:
	go mod download

tidy:
	go mod tidy

clean:
	rm -f $(BINARY_PATH)
	rm -f $(TMP_BINARY)
	rm -rf bin/
	rm -rf tmp/

install:
	go install $(MAIN_PATH)

build-all:
	@mkdir -p bin
	GOOS=linux GOARCH=amd64 go build -o bin/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=amd64 go build -o bin/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=arm64 go build -o bin/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	GOOS=windows GOARCH=amd64 go build -o bin/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)