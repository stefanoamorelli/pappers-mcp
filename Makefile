.PHONY: build test test-integration conformance clean

BINARY=pappers-mcp
BUILD_DIR=.

build:
	go build -o $(BUILD_DIR)/$(BINARY) ./cmd/pappers-mcp

test:
	go test ./... -count=1

test-integration:
	go test ./... -count=1 -run Integration

conformance:
	go test ./internal/tools/ -count=1 -run TestConformance

clean:
	rm -f $(BUILD_DIR)/$(BINARY)
