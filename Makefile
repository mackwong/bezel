
LDFLAGS=-ldflags="-X gitlab.bj.sensetime.com/diamond/service-providers/bezel/cmd.Version=$(shell git describe) -w -s -extldflags -static"
BINARY=bin/bezel

build: ## Build
	go build -o ${BINARY} ${LDFLAGS} ./cmd/bezel/main.go

test:
	echo 'mode: atomic' > coverage.txt && go test -covermode=atomic -coverprofile=coverage.txt -v -run="Test*" -timeout=30s ./...
	go tool cover -html=coverage.txt -o coverage.html

linux: ## Build an arm64/amd64 linux static binary
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 \
        go build -o ${BINARY}-linux ${LDFLAGS} cmd/bezel/main.go

darwin: ## Build a macos static binary
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 \
        go build -o ${BINARY}-darwin ${LDFLAGS} cmd/bezel/main.go

windows: ## Build a windows static binary
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 \
        go build -o ${BINARY}-windows.exe ${LDFLAGS} cmd/bezel/main.go

clean: ## Cleans up build artifacts
	rm -rf bin/ coverage.txt

lint: ## Run all the linters
	golangci-lint run --fast --deadline 3m  --skip-dirs vendor ./...

all: clean linux darwin windows
	tar cfz bezel-$(shell git describe).tar.gz bin/

.PHONY: build test darwin linux windows lint clean all
