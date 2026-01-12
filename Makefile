APP_MCP=mcp-github
BIN_DIR=bin

VERSION?=dev
COMMIT?=$(shell git rev-parse --short HEAD 2>/dev/null || echo none)
BUILT?=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

LDFLAGS=-ldflags "\
  -X 'github.com/aidashovv/pr-review-bot/internal/buildinfo.Version=$(VERSION)' \
  -X 'github.com/aidashovv/pr-review-bot/internal/buildinfo.Commit=$(COMMIT)' \
  -X 'github.com/aidashovv/pr-review-bot/internal/buildinfo.Built=$(BUILT)' \
"

.PHONY: build build-mcp run-mcp

build: build-mcp

build-mcp:
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=0 go build $(LDFLAGS) -o $(BIN_DIR)/$(APP_MCP) ./cmd/mcp-github

run-mcp:
	go run ./cmd/mcp-github
