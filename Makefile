# -----------------------------
# Warning: 100% AI-generated
# -----------------------------

ROOT := $(CURDIR)
DST := $(ROOT)/build

PNPM := pnpm
GO := go
CURL := curl
UNZIP := unzip

FRONTEND_ROOT := $(ROOT)/frontend
BACKEND_ROOT := $(ROOT)/backend

VERSION ?= dev
COMMIT ?= unknown-commit
BUILD_TIME ?= unknown-time

XRAY_VERSION ?= v26.5.9

# Целевая платформа сборки (можно переопределить: make build GOOS=linux GOARCH=arm64)
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
CGO_ENABLED ?= 0

# -----------------------------
# DEFAULT
# -----------------------------

.PHONY: all build
all: build

# Собирает и node, и nodeman
build: build_node build_nodeman

# -----------------------------
# CLEAN
# -----------------------------

.PHONY: clean clean_node clean_nodeman clean_frontend clean_xray

# Общий clean — чистит всё
clean: clean_node clean_nodeman clean_frontend
	rm -rf $(DST)

# Чистит только node (включая скачанный xray и его data/)
clean_node:
	rm -rf $(DST)/xray-node

# Чистит только nodeman
clean_nodeman:
	rm -rf $(DST)/xray-nodeman

clean_frontend:
	rm -rf $(FRONTEND_ROOT)/admpage/dist
	rm -rf $(FRONTEND_ROOT)/userpage/dist

# -----------------------------
# FRONTEND (нужен только для nodeman)
# -----------------------------

.PHONY: deps_frontend gen_frontend build_frontend

deps_frontend:
	cd $(FRONTEND_ROOT) && $(PNPM) install

gen_frontend: deps_frontend
	cd $(FRONTEND_ROOT) && $(PNPM) run gen

build_frontend: gen_frontend
	@echo "==> Building frontend"
	cd $(FRONTEND_ROOT) && $(PNPM) run build

# -----------------------------
# EMBED FRONTEND INTO NODEMAN
# -----------------------------

.PHONY: embed_frontend

embed_frontend: build_frontend
	@echo "==> Embedding frontend into nodeman"

	rm -rf $(BACKEND_ROOT)/nodeman/internal/pages/admpage
	rm -rf $(BACKEND_ROOT)/nodeman/internal/pages/userpage

	cp -rp $(FRONTEND_ROOT)/admpage/dist \
		$(BACKEND_ROOT)/nodeman/internal/pages/admpage

	cp -rp $(FRONTEND_ROOT)/userpage/dist \
		$(BACKEND_ROOT)/nodeman/internal/pages/userpage

# -----------------------------
# TOOLS
# -----------------------------

.PHONY: tools

GO_TOOLS := \
	github.com/ogen-go/ogen/cmd/ogen@latest \
	github.com/jmattheis/goverter/cmd/goverter@latest \
	go.uber.org/mock/mockgen@latest \
	github.com/atombender/go-jsonschema@latest \
	github.com/sqlc-dev/sqlc/cmd/sqlc@latest

tools:
	@echo "==> Installing Go tools"
	@for tool in $(GO_TOOLS); do \
		echo "-> $$tool"; \
		go install $$tool; \
	done

# -----------------------------
# BACKEND: NODE
# -----------------------------

.PHONY: deps_node gen_node build_node


NODE_VERSION_PKG := github.com/XRay-Addons/xrayman/node/internal/version
NODE_LDFLAGS := \
	-X $(NODE_VERSION_PKG).Version=$(VERSION) \
	-X $(NODE_VERSION_PKG).Commit=$(COMMIT) \
	-X $(NODE_VERSION_PKG).BuildTime=$(BUILD_TIME)

deps_node:
	cd $(BACKEND_ROOT)/node && $(GO) mod download

gen_node: deps_node
	cd $(BACKEND_ROOT)/node && $(GO) generate ./...

# node зависит от xray-рантайма, но НЕ от фронтенда
build_node: gen_node download_xray
	@echo "==> Building node ($(GOOS)/$(GOARCH))"

	mkdir -p $(DST)/xray-node

	cd $(BACKEND_ROOT)/node && \
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) \
	$(GO) build -ldflags "$(NODE_LDFLAGS)" -o $(DST)/xray-node/xray-node ./cmd/main.go

# -----------------------------
# BACKEND: NODEMAN
# -----------------------------

.PHONY: deps_nodeman gen_nodeman build_nodeman

NODEMAN_VERSION_PKG := github.com/XRay-Addons/xrayman/nodeman/internal/version
NODEMAN_LDFLAGS := \
	-X $(NODEMAN_VERSION_PKG).Version=$(VERSION) \
	-X $(NODEMAN_VERSION_PKG).Commit=$(COMMIT) \
	-X $(NODEMAN_VERSION_PKG).BuildTime=$(BUILD_TIME)

deps_nodeman:
	cd $(BACKEND_ROOT)/nodeman && $(GO) mod download

gen_nodeman: deps_nodeman
	cd $(BACKEND_ROOT)/nodeman && $(GO) generate ./...

# nodeman зависит от фронтенда (админка/юзерпейдж встраиваются внутрь)
build_nodeman: gen_nodeman embed_frontend
	@echo "==> Building nodeman ($(GOOS)/$(GOARCH))"

	mkdir -p $(DST)/xray-nodeman

	cd $(BACKEND_ROOT)/nodeman && \
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) \
	$(GO) build -ldflags "$(NODEMAN_LDFLAGS)" -o $(DST)/xray-nodeman/xray-nodeman ./cmd/main.go

# -----------------------------
# XRAY DOWNLOAD
# -----------------------------

.PHONY: download_xray

XRAY_ASSET :=

ifeq ($(GOOS),darwin)
	ifeq ($(GOARCH),arm64)
		XRAY_ASSET := Xray-macos-arm64-v8a.zip
	endif
	ifeq ($(GOARCH),amd64)
		XRAY_ASSET := Xray-macos-64.zip
	endif
endif

ifeq ($(GOOS),linux)
	ifeq ($(GOARCH),amd64)
		XRAY_ASSET := Xray-linux-64.zip
	endif
	ifeq ($(GOARCH),arm64)
		XRAY_ASSET := Xray-linux-arm64-v8a.zip
	endif
endif

ifeq ($(XRAY_ASSET),)
$(error Unsupported platform: $(GOOS)/$(GOARCH))
endif

XRAY_URL := https://github.com/XTLS/Xray-core/releases/download/$(XRAY_VERSION)/$(XRAY_ASSET)

download_xray:
	@echo "==> Downloading Xray: $(XRAY_ASSET) ($(GOOS)/$(GOARCH))"
	mkdir -p $(DST)/xray-node/data

	$(CURL) -L -o $(DST)/xray.zip $(XRAY_URL)
	$(UNZIP) -j -o $(DST)/xray.zip 'xray' -d $(DST)/xray-node
	$(UNZIP) -j -o $(DST)/xray.zip 'geoip.dat' 'geosite.dat' -d $(DST)/xray-node/data
	rm -f $(DST)/xray.zip

	@echo "==> Xray ready at $(DST)/xray-node"