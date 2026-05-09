SHELL := /bin/bash

BIN_DIR := $(CURDIR)/bin
NPM_ROOT := $(CURDIR)/nodeman/web

USERPAGE_WEB_SRC := $(CURDIR)/nodeman/web/pages/userpage
ADMPAGE_WEB_SRC := $(CURDIR)/nodeman/web/pages/admpage

USERPAGE_WEB_DST := $(CURDIR)/nodeman/internal/pages/userpage
ADMPAGE_WEB_DST := $(CURDIR)/nodeman/internal/pages/admpage

NODE_SRC := $(CURDIR)/node
NODE_DST := $(BIN_DIR)/node

NODEMAN_SRC := $(CURDIR)/nodeman
NODEMAN_DST := $(BIN_DIR)/nodeman


# =========================
# FRONTEND
# =========================

.PHONY: npm-install
npm-install:
	@echo "Installing pnpm dependencies..."
	cd $(NPM_ROOT) && pnpm install


.PHONY: userpage
userpage:
	@echo "Building userpage..."
	cd $(NPM_ROOT) && pnpm run build:user
	rm -rf $(USERPAGE_WEB_DST)
	mkdir -p $(USERPAGE_WEB_DST)
	cp -r $(USERPAGE_WEB_SRC)/dist/* $(USERPAGE_WEB_DST)/


.PHONY: admpage
admpage:
	@echo "Building admpage..."
	cd $(NPM_ROOT) && pnpm run build:admin
	rm -rf $(ADMPAGE_WEB_DST)
	mkdir -p $(ADMPAGE_WEB_DST)
	cp -r $(ADMPAGE_WEB_SRC)/dist/* $(ADMPAGE_WEB_DST)/


# =========================
# BACKEND (ADDED)
# =========================

GO_TOOLS := \
	github.com/ogen-go/ogen/cmd/ogen@latest \
	github.com/jmattheis/goverter/cmd/goverter@latest \
	go.uber.org/mock/mockgen@latest \

.PHONY: tools
tools:
	@echo "Installing Go tools..."
	@for tool in $(GO_TOOLS); do \
		echo "-> $$tool"; \
		go install $$tool; \
	done

.PHONY: node_generate
node_generate: tools
	@echo "Running go generate..."
	go generate ./node/...

.PHONY: nodeman_generate
nodeman_generate: tools
	@echo "Running go generate..."
	go generate ./nodeman/...

.PHONY: node
node: node_generate
	@echo "Building node..."
	mkdir -p $(BIN_DIR)
	go build -o $(NODE_DST) $(NODE_SRC)/cmd


.PHONY: nodeman
nodeman: node nodeman_generate
	@echo "Building nodeman..."
	mkdir -p $(BIN_DIR)
	go build -o $(NODEMAN_DST) $(NODEMAN_SRC)/cmd


# =========================
# FULL BUILD
# =========================

.PHONY: build
build: npm-install userpage admpage node_generate nodeman_generate node nodeman
	@echo "Full build done (frontend + backend)."