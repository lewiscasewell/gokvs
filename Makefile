APP_NAME := mini-redis
VERSION := 1.0.0
BUILD_DIR := dist

PLATFORMS := \
  linux/amd64 \
  linux/arm64 \
  darwin/amd64 \
  darwin/arm64 \
  windows/amd64

.PHONY: all build run clean

all: build

build:
	@echo "🔨 Building $(APP_NAME) for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	@for platform in $(PLATFORMS); do \
		GOOS=$${platform%%/*} \
		GOARCH=$${platform##*/} \
		OUTPUT=$(BUILD_DIR)/$(APP_NAME)-$${platform%%/*}-$${platform##*/}; \
		if [ "$${platform%%/*}" = "windows" ]; then OUTPUT=$$OUTPUT.exe; fi; \
		echo "⚙️  Building $$OUTPUT..."; \
		GOOS=$${platform%%/*} GOARCH=$${platform##*/} go build -o $$OUTPUT server/*.go || exit 1; \
	done
	@echo "✅ Build complete."

run:
	@echo "🚀 Running $(APP_NAME)..."
	go run server/*.go

clean:
	@echo "🧹 Cleaning up..."
	@rm -rf $(BUILD_DIR)
	@echo "✅ Clean complete."
