BINARY     := stereoctl
MODULE     := $(shell go list -m)
VERSION    := $(shell grep 'var version' cmd/root.go | sed 's/.*"\(.*\)"/\1/')
LDFLAGS    := -ldflags "-s -w -X $(MODULE)/cmd.version=$(VERSION)"
BIN_DIR    := bin
BUILD_DIR  := dist

GOOS   ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

.DEFAULT_GOAL := help

# ── Help ──────────────────────────────────────────────────────────────────────
.PHONY: help
help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "  build         Build binary for the current platform ($(GOOS)/$(GOARCH))"
	@echo "  build-all     Cross-compile for linux/amd64, darwin/amd64, darwin/arm64, windows/amd64"
	@echo "  run           Run the tool (pass ARGS=\"convert file.mkv\")"
	@echo "  test          Run tests"
	@echo "  lint          Run go vet"
	@echo "  fmt           Format source code"
	@echo "  tidy          Tidy and verify go.mod / go.sum"
	@echo "  clean         Remove build artifacts"
	@echo "  version       Show current version"

# ── Build ─────────────────────────────────────────────────────────────────────
.PHONY: build
build: $(BIN_DIR)
	CGO_ENABLED=0 go build $(LDFLAGS) -o $(BIN_DIR)/$(BINARY) .

$(BIN_DIR):
	mkdir -p $(BIN_DIR)

.PHONY: build-all
build-all: $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY)-linux-amd64 .
	CGO_ENABLED=0 GOOS=darwin  GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY)-darwin-amd64 .
	CGO_ENABLED=0 GOOS=darwin  GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY)-darwin-arm64 .
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY)-windows-amd64.exe .
	@echo "Binaries available in $(BUILD_DIR)/"

$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

# ── Run ───────────────────────────────────────────────────────────────────────
.PHONY: run
run:
	go run . $(ARGS)

# ── Quality ───────────────────────────────────────────────────────────────────
.PHONY: test
test:
	go test ./...

.PHONY: lint
lint:
	go vet ./...

.PHONY: fmt
fmt:
	gofmt -w .

.PHONY: tidy
tidy:
	go mod tidy
	go mod verify

# ── Util ──────────────────────────────────────────────────────────────────────
.PHONY: clean
clean:
	rm -rf $(BIN_DIR) $(BUILD_DIR)

.PHONY: version
version:
	@echo "v$(VERSION)"
