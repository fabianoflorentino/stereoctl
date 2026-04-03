BINARY     := stereoctl
MODULE     := $(shell go list -m)
VERSION    := $(shell cat VERSION 2>/dev/null || git describe --tags --abbrev=0 2>/dev/null | sed 's/^v//' || echo "dev")
LDFLAGS    := -ldflags "-s -w -X $(MODULE)/cmd.version=$(VERSION)"
BIN_DIR    := bin
BUILD_DIR  := dist
DIST_OUTPUT := $(BIN_DIR)/dist

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
	@echo "  dist          Package cross-built artifacts into $(DIST_OUTPUT)"
	@echo "  release       Build and package release artifacts into $(DIST_OUTPUT)"
	@echo "  hooks-install Install local git hooks (lefthook)"
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
build-all:
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY)-linux-amd64 .
	CGO_ENABLED=0 GOOS=darwin  GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY)-darwin-amd64 .
	CGO_ENABLED=0 GOOS=darwin  GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY)-darwin-arm64 .
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY)-windows-amd64.exe .
	@echo "Binaries available in $(BUILD_DIR)/"

.PHONY: dist
dist: build-all $(DIST_OUTPUT)
	@echo "Packaging artifacts into $(DIST_OUTPUT)"
	@for f in $(BUILD_DIR)/* ; do \
		name=$$(basename $$f) ; \
		case $$name in \
			*.exe) zip -j $(DIST_OUTPUT)/$$name.zip $$f ;; \
			*) tar -czf $(DIST_OUTPUT)/$$name.tar.gz -C $(BUILD_DIR) $$name ;; \
		esac ; \
	done
	@echo "Packaged artifacts in $(DIST_OUTPUT)"

.PHONY: hooks-install
hooks-install:
	@echo "Installing lefthook hooks via script..."
	@bash scripts/install-lefthook.sh

.PHONY: release
release: dist
	@echo "Release artifacts ready in $(DIST_OUTPUT)"

$(DIST_OUTPUT):
	mkdir -p $(DIST_OUTPUT)

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
