GO_CMD = go
GO_BUILD_CMD = $(GO_CMD) build
APP_NAME = go-tms

RELEASE_DIR = ./release

LDFLAGS = -ldflags="-s -w"

.PHONY: all clean linux darwin

all: linux darwin

linux:
	@echo "Building for Linux (x64_86)..."
	mkdir -p $(RELEASE_DIR)/x64_86
	GOOS=linux GOARCH=amd64 $(GO_BUILD_CMD) $(LDFLAGS) -o $(RELEASE_DIR)/x64_86/$(APP_NAME)
	@echo "Linux build complete."

darwin:
	@echo "Building for macOS (aarch64)..."
	mkdir -p $(RELEASE_DIR)/darwin
	GOOS=darwin GOARCH=arm64 $(GO_BUILD_CMD) $(LDFLAGS) -o $(RELEASE_DIR)/darwin/$(APP_NAME)
	@echo "macOS build complete."


clean:
	@echo "Cleaning up..."
	rm -rf $(RELEASE_DIR)
	@echo "Cleanup complete."

