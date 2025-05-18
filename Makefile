APP_NAME := orbital
BUILD_DIR := build

GOARCHES := amd64 arm64
GOOS := linux

.PHONY: all clean build

all: clean build

build:
	@mkdir -p $(BUILD_DIR)
	@for arch in $(GOARCHES); do \
		echo "Building for $(GOOS)/$$arch..."; \
		GOOS=$(GOOS) GOARCH=$$arch CGO_ENABLED=0 go build -o $(BUILD_DIR)/$(APP_NAME)-$(GOOS)-$$arch .; \
	done

clean:
	rm -rf $(BUILD_DIR)
