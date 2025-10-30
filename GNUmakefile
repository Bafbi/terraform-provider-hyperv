default: build

# Build and install the provider locally
.PHONY: build
build:
	go build -o terraform-provider-hyperv
	go install .

# Install the provider to GOBIN for local development
.PHONY: install
install: build
	@echo "Provider installed to $(shell go env GOBIN)/terraform-provider-hyperv"
	@echo ""
	@echo "To use with Terraform, set:"
	@echo "  export TF_CLI_CONFIG_FILE=$(PWD)/.terraformrc"

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

# Run unit tests only
.PHONY: test
test:
	go test ./... -v -short $(TESTARGS) -timeout 10m

# Format code
.PHONY: fmt
fmt:
	go fmt ./...
	terraform fmt -recursive ./examples/

# Lint code
.PHONY: lint
lint:
	golangci-lint run

# Generate documentation
.PHONY: docs
docs:
	go generate ./...

# Clean build artifacts
.PHONY: clean
clean:
	rm -f terraform-provider-hyperv
	go clean

# Initialize a test directory with the local provider
.PHONY: test-init
test-init: install
	@echo "Initializing test environment..."
	@cd examples/vm-from-scratch && \
		rm -rf .terraform .terraform.lock.hcl && \
		TF_CLI_CONFIG_FILE=$(PWD)/.terraformrc terraform init

# Run terraform plan in test directory
.PHONY: test-plan
test-plan: install
	@cd examples/vm-from-scratch && \
		TF_CLI_CONFIG_FILE=$(PWD)/.terraformrc terraform plan

# Run terraform apply in test directory
.PHONY: test-apply
test-apply: install
	@cd examples/vm-from-scratch && \
		TF_CLI_CONFIG_FILE=$(PWD)/.terraformrc terraform apply

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build       - Build and install the provider"
	@echo "  install     - Install the provider to GOBIN for local development"
	@echo "  test        - Run unit tests"
	@echo "  testacc     - Run acceptance tests"
	@echo "  fmt         - Format Go and Terraform code"
	@echo "  lint        - Run linter"
	@echo "  docs        - Generate documentation"
	@echo "  clean       - Clean build artifacts"
	@echo "  test-init   - Initialize test directory with local provider"
	@echo "  test-plan   - Run terraform plan in test directory"
	@echo "  test-apply  - Run terraform apply in test directory"
	@echo "  help        - Show this help message"
