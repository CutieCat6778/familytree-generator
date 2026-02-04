.PHONY: build run clean test lint install

# Binary name
BINARY=familytree

# Build the CLI binary
build:
	go build -o $(BINARY) ./cmd/familytree

# Build the server binary
build-server:
	go build -o $(BINARY)-server ./cmd/server

# Build both binaries
build-all-local: build build-server

# Run the application
run: build
	./$(BINARY)

# Run with verbose output
run-verbose: build
	./$(BINARY) -verbose

# List available countries
list-countries: build
	./$(BINARY) -list-countries

# Generate a sample tree
sample: build
	./$(BINARY) -country united-states -generations 3 -seed 42 -verbose

# Generate JSON output
sample-json: build
	./$(BINARY) -country japan -generations 3 -format json -output tree.json -verbose

# Clean build artifacts
clean:
	rm -f $(BINARY)
	rm -f *.csv *.json
	rm -f family_tree*

# Run tests
test:
	go test -v ./...

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Install dependencies
deps:
	go mod tidy

# Build for all platforms
build-all:
	GOOS=linux GOARCH=amd64 go build -o $(BINARY)-linux-amd64 ./cmd/familytree
	GOOS=darwin GOARCH=amd64 go build -o $(BINARY)-darwin-amd64 ./cmd/familytree
	GOOS=darwin GOARCH=arm64 go build -o $(BINARY)-darwin-arm64 ./cmd/familytree
	GOOS=windows GOARCH=amd64 go build -o $(BINARY)-windows-amd64.exe ./cmd/familytree

# Server
server: build-server
	./$(BINARY)-server -port 8080

# Web development
web-install:
	cd web && npm install

web-dev:
	cd web && npm run dev

web-build:
	cd web && npm run build

# Full stack development
dev: build-server web-build
	./$(BINARY)-server -port 8080 -web ./web/dist

# Help
help:
	@echo "Available targets:"
	@echo ""
	@echo "Go CLI:"
	@echo "  build         - Build the CLI binary"
	@echo "  build-server  - Build the API server"
	@echo "  run           - Build and run CLI"
	@echo "  run-verbose   - Build and run CLI with verbose output"
	@echo "  list-countries- List available countries"
	@echo "  sample        - Generate a sample tree"
	@echo "  sample-json   - Generate a sample JSON tree"
	@echo "  clean         - Remove build artifacts"
	@echo "  test          - Run tests"
	@echo "  fmt           - Format code"
	@echo "  lint          - Lint code"
	@echo "  deps          - Install dependencies"
	@echo "  build-all     - Build for all platforms"
	@echo ""
	@echo "Server:"
	@echo "  server        - Start the API server on port 8080"
	@echo "  dev           - Build and start full stack (server + web)"
	@echo ""
	@echo "Web Visualization:"
	@echo "  web-install   - Install web dependencies"
	@echo "  web-dev       - Start web dev server"
	@echo "  web-build     - Build web for production"
