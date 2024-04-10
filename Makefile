OS_NAME := $(shell uname)
ifeq ($(OS_NAME), Darwin)
OPEN := open
else
OPEN := xdg-open
endif

BINARY_NAME=get-next-version
MODULE_PATH=github.com/thenativeweb/${BINARY_NAME}

COMMON_FLAGS=-ldflags="-X '$(MODULE_PATH)/appversion.AppVersion=$(VERSION)'"

BUILD_DIR=./build

qa: analyze test

analyze:
	@go vet ./...

test:
	@go test -cover ./...

benchmark:
	@go test -bench=. -run=^$$ ./...

coverage:
	@mkdir -p ./coverage
	@go test -coverprofile=./coverage/cover.out ./...
	@go tool cover -html=./coverage/cover.out -o ./coverage/cover.html
	@$(OPEN) ./coverage/cover.html

clean:
	@rm -rf build/

detect-version:
	$(eval VERSION=$(shell git tag --points-at HEAD))
	$(eval VERSION=$(or $(VERSION), (version unavailable)))

	@echo "Current Version: ${VERSION}"

build-auto: qa clean detect-version
	@go build $(COMMON_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)

build-darwin-amd64: qa clean detect-version
	@GOOS=darwin GOARCH=amd64 go build $(COMMON_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64

build-darwin-arm64: qa clean detect-version
	@GOOS=darwin GOARCH=arm64 go build $(COMMON_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64

build-linux-amd64: qa clean detect-version
	@GOOS=linux GOARCH=amd64 go build $(COMMON_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64

build-windows-amd64: qa clean detect-version
	@GOOS=windows GOARCH=amd64 go build $(COMMON_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe

build: build-auto

build-all: build-darwin-amd64 build-darwin-arm64 build-linux-amd64 build-windows-amd64

.PHONY: analyze \
		benchmark \
		build \
		build-all \
		build-auto \
		build-darwin-amd64 \
		build-darwin-arm64 \
		build-linux-amd64 \
		build-windows-amd64 \
		clean \
		coverage \
		detect-version \
		qa \
		test
