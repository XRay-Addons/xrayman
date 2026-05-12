# -----------------------------
# Warning: 100% AI-generated
# -----------------------------

ROOT := $(CURDIR)
DST := $(ROOT)/build

PNPM := pnpm
GO := go

FRONTEND_ROOT := $(ROOT)/frontend
BACKEND_ROOT := $(ROOT)/backend

# -----------------------------
# ALL FILES (recursive tracking)
# -----------------------------

FRONTEND_FILES := $(shell find $(FRONTEND_ROOT) -type f)
BACKEND_NODE_FILES := $(shell find $(BACKEND_ROOT)/node -type f)
BACKEND_NODMAN_FILES := $(shell find $(BACKEND_ROOT)/nodeman -type f)

.PHONY: all build clean

all: build

build: build_backend

# -----------------------------
# CLEAN
# -----------------------------

clean:
	rm -rf $(DST)

clean_frontend:
	rm -rf $(FRONTEND_ROOT)/admpage/dist
	rm -rf $(FRONTEND_ROOT)/userpage/dist

clean_backend:
	rm -rf $(DST)

# -----------------------------
# FRONTEND
# -----------------------------

.PHONY: build_frontend gen_frontend deps_frontend embed_frontend

deps_frontend:
	cd $(FRONTEND_ROOT) && $(PNPM) install

gen_frontend: deps_frontend
	cd $(FRONTEND_ROOT) && $(PNPM) run gen

# rebuild only if ANY frontend file changed (recursive)
build_frontend: $(FRONTEND_FILES)
	@echo "Building frontend..."
	cd $(FRONTEND_ROOT) && $(PNPM) run build

# embed depends on real frontend build
embed_frontend: build_frontend
	@echo "Embedding frontend into backend..."

	rm -rf $(BACKEND_ROOT)/nodeman/internal/pages/admpage
	rm -rf $(BACKEND_ROOT)/nodeman/internal/pages/userpage

	cp -rp $(FRONTEND_ROOT)/admpage/dist \
		$(BACKEND_ROOT)/nodeman/internal/pages/admpage

	cp -rp $(FRONTEND_ROOT)/userpage/dist \
		$(BACKEND_ROOT)/nodeman/internal/pages/userpage

# -----------------------------
# BACKEND
# -----------------------------

GO_TOOLS := \
	github.com/ogen-go/ogen/cmd/ogen@latest \
	github.com/jmattheis/goverter/cmd/goverter@latest \
	go.uber.org/mock/mockgen@latest

tools:
	@echo "Installing Go tools..."
	@for tool in $(GO_TOOLS); do \
		echo "-> $$tool"; \
		go install $$tool; \
	done

deps_backend:
	cd $(BACKEND_ROOT)/node && $(GO) mod download
	cd $(BACKEND_ROOT)/nodeman && $(GO) mod download

gen_backend: deps_backend
	cd $(BACKEND_ROOT)/node && $(GO) generate ./...
	cd $(BACKEND_ROOT)/nodeman && $(GO) generate ./...

# -----------------------------
# BUILD BACKEND
# -----------------------------

build_backend: gen_backend embed_frontend
	@echo "Building backend..."

	mkdir -p $(DST)/xray-node
	mkdir -p $(DST)/xray-nodeman

	cd $(BACKEND_ROOT)/node && \
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) \
	$(GO) build -o $(DST)/xray-node/xray-node ./cmd/main.go

	cd $(BACKEND_ROOT)/nodeman && \
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) \
	$(GO) build -o $(DST)/xray-nodeman/xray-nodeman ./cmd/main.go


# -----------------------------
#  DOWNLOAD XRAY TOOLS
# -----------------------------

# # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #
# XRAY DOWNLOAD
# # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #

XRAY_VERSION ?= v26.5.9
XRAY_DST := $(DST)/xray

GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)

XRAY_ASSET :=

# detect version
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

.PHONY: xray clean_xray

xray:
	@echo "==> Downloading Xray: $(XRAY_ASSET)"
	rm -rf $(XRAY_DST)
	mkdir -p $(XRAY_DST)
	curl -L -o $(DST)/xray.zip $(XRAY_URL)
	
	@echo "unzip Xray: $(XRAY_DST)"
	unzip -o $(DST)/xray.zip -d $(XRAY_DST)
	rm -f $(DST)/xray.zip
	@echo "copy xray geodata to $(DST)/xray-node/xray-node/data/"
	mkdir -p $(DST)/xray-node/data/
	mv -r $(XRAY_DST)/geoip.dat $(DST)/xray-node/data/
	mv -r $(XRAY_DST)/geosite.dat $(DST)/xray-node/data/
	@echo "==> Xray ready at $(XRAY_DST)"

clean_xray:
	rm -rf $(XRAY_DST)