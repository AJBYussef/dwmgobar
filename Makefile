# Define the binary name and module name
BINARY_NAME := dwmgobar

# Directories
SRC_DIR := .
DEST_DIR := $(HOME)/.local/bin
BUILD_DIR := ./build

# Find all Go source files
SRC_FILES := $(wildcard $(SRC_DIR)/*.go)

# Build the application
build: $(BUILD_DIR)/$(BINARY_NAME)

$(BUILD_DIR)/$(BINARY_NAME): $(SRC_FILES)
	@echo "Building $(BINARY_NAME)..."
	@go build -o $@ $(SRC_DIR)

# Run the application
run: build
	@echo "Running $(BINARY_NAME)..."
	@./$(BUILD_DIR)/$(BINARY_NAME) 

install: build
	@echo "Installing $(BINARY_NAME)..."
	mkdir -p $(DEST_DIR)
	cp -f $(BUILD_DIR)/$(BINARY_NAME) $(DEST_DIR)
	chmod 755 $(DEST_DIR)/$(BINARY_NAME)
# Clean the build artifacts
clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)
.PHONY: build run clean
