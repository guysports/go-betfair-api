# Go project helpers

GO ?= go
GOPATH ?= $(shell $(GO) env GOPATH)
ifeq ($(strip $(GOPATH)),)
GOPATH := $(shell $(GO) env GOPATH)/bin
endif
GOBIN ?= $(shell $(GO) env GOBIN)
ifeq ($(strip $(GOBIN)),)
GOBIN := $(shell $(GO) env GOPATH)/bin
endif

GOLANGCI_LINT_VERSION ?= v2.12.2
GOFUMPT_VERSION ?= v0.11.0

.PHONY: lint vet test build-windows build-darwin install-tools install-golangci-lint install-gofumpt

install-tools: install-golangci-lint install-gofumpt

install-golangci-lint:
	curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b $(GOPATH)/bin v2.12.2

install-gofumpt:
	@if [ -x $(GOBIN)/gofumpt ]; then echo "gofumpt already installed"; else mkdir -p $(GOBIN) && $(GO) install mvdan.cc/gofumpt@$(GOFUMPT_VERSION); fi

lint:
	@command -v $(GOBIN)/golangci-lint >/dev/null 2>&1 || { echo "golangci-lint is not installed; run 'make install-golangci-lint'"; exit 1; }
	$(GOBIN)/golangci-lint run

vet:
	go vet ./...

test:
	go test ./...

build-windows:
	mkdir -p bin
	GOOS=windows GOARCH=amd64 go build -o bin/go-betfair-api-windows-amd64.exe .

build-darwin:
	mkdir -p bin
	GOOS=darwin GOARCH=amd64 go build -o bin/go-betfair-api-darwin-amd64 .
