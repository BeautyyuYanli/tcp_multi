.PHONY: build run test clean lint build-arm64

# Build directory
BUILD_DIR=./bin

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOLINT=golangci-lint

# Get the program name from the second argument
PROGRAM=$(word 2,$(MAKECMDGOALS))
ifeq ($(PROGRAM),)
PROGRAM_ERROR=true
endif

# Build the application
build:
ifeq ($(PROGRAM_ERROR),true)
	@echo "Usage: make build <program_name>"
	@exit 1
else
	mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(PROGRAM) ./cmd/$(PROGRAM)
endif

# Build for Linux ARM64 (ARMv8)
build-arm64:
ifeq ($(PROGRAM_ERROR),true)
	@echo "Usage: make build-arm64 <program_name>"
	@exit 1
else
	mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=arm64 $(GOBUILD) -o $(BUILD_DIR)/$(PROGRAM)-linux-arm64 ./cmd/$(PROGRAM)
	@echo "Built for Linux ARM64: $(BUILD_DIR)/$(PROGRAM)-linux-arm64"
endif

# Run the application
run:
ifeq ($(PROGRAM_ERROR),true)
	@echo "Usage: make run <program_name>"
	@exit 1
else
	$(GORUN) ./cmd/$(PROGRAM)
endif

# Test the application
test:
	$(GOTEST) ./...

# Clean build artifacts
clean:
	rm -rf $(BUILD_DIR)

# Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Run linter
lint:
	$(GOLINT) run

# Install golangci-lint if not installed
install-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

# Default target
all: clean

# Handle arbitrary targets (to catch the program name)
%:
	@: 