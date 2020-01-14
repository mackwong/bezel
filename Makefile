build: ## Build
	go build -o bin/bezel cmd/bezel/main.go

test: ## Run the tests
	go test -v $(shell go list ./... | grep -v /vendor/ | grep -v /test/)

arm64: ## Build an arm64 static binary
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 \
        go build -v -o bin/bezel-arm64 -ldflags="-w -s -extldflags -static" cmd/bezel/main.go

amd64: ## Build an amd64 static binary
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
        go build -v -o bin/bezel-amd64 -ldflags="-w -s -extldflags -static" cmd/bezel/main.go

windows: ## Build a windows static binary
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 \
        go build -v -o bin/bezel-windows.exe -ldflags="-w -s -extldflags -static" cmd/bezel/main.go

clean: ## Cleans up build artifacts
	rm -f bezel

lint: ## Run all the linters
	golangci-lint run --fast --deadline 3m  --skip-dirs vendor ./...

.PHONY: build test arm64 amd64 windows lint clean
