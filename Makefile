.PHONY: help hooks test build clean

# Default target
help: ## Show this help message
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*##"; printf "\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  %-15s %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

hooks: ## Install Git hooks for pre-commit and pre-push validation
	@./install-git-hooks.sh

test: ## Run all tests
	@go test -v -race ./internal/...

test-short: ## Run fast tests only
	@go test -short -race ./internal/...

build: ## Build the anvil binary
	@go build -o anvil cmd/anvil/main.go

clean: ## Clean build artifacts
	@rm -f anvil anvil-test* coverage.out

lint: ## Run linter (requires golangci-lint)
	@golangci-lint run

fmt: ## Format Go code
	@gofmt -w .

pre-push: ## Run the full pre-push test suite
	@./pre-push-test.sh

deps: ## Download and verify dependencies
	@go mod download && go mod verify

# Install hooks automatically when downloading dependencies
deps-with-hooks: deps hooks ## Download dependencies and install hooks
