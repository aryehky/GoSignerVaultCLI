.PHONY: build test clean lint

# Build the CLI
build:
	go build -o bin/gosignervaultcli main.go

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Run linter
lint:
	golangci-lint run

# Install dependencies
deps:
	go mod tidy

# Create a new release
release:
	@echo "Creating release..."
	@version=$$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.1.0"); \
	git tag -a $$version -m "Release $$version"; \
	git push origin $$version

# Install the CLI
install: build
	cp bin/gosignervaultcli /usr/local/bin/

# Run the CLI
run: build
	./bin/gosignervaultcli

# Help
help:
	@echo "Available targets:"
	@echo "  build    - Build the CLI"
	@echo "  test     - Run tests"
	@echo "  clean    - Clean build artifacts"
	@echo "  lint     - Run linter"
	@echo "  deps     - Install dependencies"
	@echo "  release  - Create a new release"
	@echo "  install  - Install the CLI"
	@echo "  run      - Run the CLI"
	@echo "  help     - Show this help message" 