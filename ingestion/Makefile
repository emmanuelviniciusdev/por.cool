.PHONY: build run test test-coverage clean docker-build docker-run deps lint

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=ingestion

# Build the application
build:
	CGO_ENABLED=0 $(GOBUILD) -ldflags="-w -s" -o $(BINARY_NAME) .

# Run the application
run:
	$(GOCMD) run main.go

# Run tests
test:
	$(GOTEST) -v ./...

# Run tests with coverage
test-coverage:
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Clean build artifacts
clean:
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html

# Build Docker image
docker-build:
	docker build -t porcool-ingestion .

# Run with Docker Compose
docker-run:
	docker-compose up -d

# Stop Docker Compose
docker-stop:
	docker-compose down

# View logs
docker-logs:
	docker-compose logs -f ingestion

# Lint the code (requires golangci-lint)
lint:
	golangci-lint run ./...

# Format code
fmt:
	$(GOCMD) fmt ./...

# Verify dependencies
verify:
	$(GOMOD) verify

# Show help
help:
	@echo "Available targets:"
	@echo "  build          - Build the binary"
	@echo "  run            - Run the application locally"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  deps           - Download and tidy dependencies"
	@echo "  clean          - Remove build artifacts"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run with Docker Compose"
	@echo "  docker-stop    - Stop Docker Compose"
	@echo "  docker-logs    - View container logs"
	@echo "  lint           - Run linter"
	@echo "  fmt            - Format code"
	@echo "  verify         - Verify dependencies"
