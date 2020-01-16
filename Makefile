build: ## Build
	go build -o bin/bezel cmd/bezel/main.go

test:
	echo 'mode: atomic' > coverage.txt && go test -covermode=atomic -coverprofile=coverage.txt -v -run="Test*" -timeout=30s ./...
	go tool cover -html=coverage.txt -o coverage.html

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
