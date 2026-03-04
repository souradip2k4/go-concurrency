.PHONY: build clean run help

# Get the absolute path of the root directory
ROOT_DIR := $(shell git rev-parse --show-toplevel)
# Get the absolute path of where the user is CURRENTLY typing the command
CURRENT_DIR := $(PWD)
# Calculate the relative path from root to current (e.g., mutexes/mutex-implementation)
RELATIVE_PATH := $(subst $(ROOT_DIR)/,,$(CURRENT_DIR))

# Variables for the build system
BINARY_NAME := main
BINARY_PATH := $(ROOT_DIR)/bin/$(RELATIVE_PATH)/$(BINARY_NAME)

# Find all main.go files and translate them into their future binary paths
MAIN_SRCS := $(shell find . -name "main.go" -not -path "*/.*" -not -path "./bin/*")
BINARIES := $(patsubst ./%/main.go,bin/%/main,$(MAIN_SRCS))

## build: Compiles all Go examples into the /bin directory (Incremental Build)
build: $(BINARIES)
	@echo "Build complete!"

# This advanced Make rule tells it that a binary depends on ALL .go files in its directory
# It will ONLY recompile if one of those Go files has changed since the last build
.SECONDEXPANSION:
bin/%/main: $$(wildcard %/*.go)
	@echo "  -> Building $@"
	@mkdir -p $(dir $@)
	@cd $* && go build -o $(ROOT_DIR)/$@ .

## clean: Deletes the /bin directory
clean:
	@echo "Cleaning up binaries..."
	@rm -rf bin/
	@echo "Done!"

## run: Runs the binary corresponding to your current directory
run:
	@if [ "$(RELATIVE_PATH)" = "" ] || [ "$(RELATIVE_PATH)" = "." ]; then \
		echo "Error: You are in the root directory. Please 'cd' into an example folder first."; \
		exit 1; \
	fi; \
	if [ -f "$(BINARY_PATH)" ]; then \
		echo "--- Running binary: $(RELATIVE_PATH)/$(BINARY_NAME) ---"; \
		"$(BINARY_PATH)"; \
	else \
		echo "Error: No compiled binary found for this path."; \
		echo "Expected at: bin/$(RELATIVE_PATH)/$(BINARY_NAME)"; \
		echo "Try running 'make build' from the root first."; \
		exit 1; \
	fi

## help: Shows available commands
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' |  sed -e 's/^/ /'
