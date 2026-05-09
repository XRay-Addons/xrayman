ROOT := $(CURDIR)
DST := $(ROOT)/build

# Утилиты
PNPM := pnpm
GO := go

# Пути
FRONTEND_ROOT := $(ROOT)/frontend
FRONTEND_DST := $(DST)/frontend
BACKEND_ROOT := $(ROOT)/backend
BACKEND_DST := $(DST)/backend

.PHONY: all build clean install

all: build

build: build_backend

clean: clean_frontend clean_backend
	rm -rf $(DST)

# # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #
# Frontend
# # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #

.PHONY: all_frontend clean_frontend gen_frontend build_frontend

all_frontend: build_frontend

deps_frontend:
	cd $(FRONTEND_ROOT) && $(PNPM) install

gen_frontend: deps_frontend
	@echo "Generating frontend..."
	cd $(FRONTEND_ROOT) && $(PNPM) run gen

build_frontend: gen_frontend
	@echo "Building frontend apps..."
	mkdir -p $(FRONTEND_DST)
	cd $(FRONTEND_ROOT) && $(PNPM) run build
	cp -rp $(FRONTEND_ROOT)/admpage/dist $(FRONTEND_DST)/admpage
	cp -rp $(FRONTEND_ROOT)/userpage/dist $(FRONTEND_DST)/userpage

clean_frontend:
	rm -rf $(FRONTEND_DST)
	rm -rf $(FRONTEND_ROOT)/admpage/dist
	rm -rf $(FRONTEND_ROOT)/userpage/dist

# # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #
# Backend
# # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #

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

.PHONY: all_backend clean_backend gen_backend embed_frontend build_backend

all_backend: build_backend

deps_backend:
	cd $(BACKEND_ROOT)/node && $(GO) mod download
	cd $(BACKEND_ROOT)/nodeman && $(GO) mod download

gen_backend: deps_backend
	@echo "Generating Backend..."
	cd $(BACKEND_ROOT)/node && $(GO) generate ./...
	cd $(BACKEND_ROOT)/nodeman && $(GO) generate ./...

embed_frontend: build_frontend
	@echo "Embedding Frontend into Backend..."
	mkdir -p $(BACKEND_ROOT)/nodeman/internal/pages
	cp -rp $(FRONTEND_DST)/admpage $(BACKEND_ROOT)/nodeman/internal/pages/
	cp -rp $(FRONTEND_DST)/userpage $(BACKEND_ROOT)/nodeman/internal/pages/

GEOIP_URL := https://raw.githubusercontent.com/Loyalsoldier/v2ray-rules-dat/release/geoip.dat
GEOSITE_URL := https://raw.githubusercontent.com/Loyalsoldier/v2ray-rules-dat/release/geosite.dat

NODE_BIN := $(BACKEND_DST)/xray-node
NODEMAN_BIN := $(BACKEND_DST)/xray-nodeman

$(NODE_BIN): gen_backend
	@echo "Building xray-node..."
	cd $(BACKEND_ROOT)/node && $(GO) build -o $@ ./cmd/main.go

$(NODEMAN_BIN): gen_backend embed_frontend
	@echo "Building xray-nodeman..."
	cd $(BACKEND_ROOT)/nodeman && $(GO) build -o $@ ./cmd/main.go

GEOIP_DAT := $(BACKEND_DST)/geoip.dat
GEOSITE_DAT := $(BACKEND_DST)/geosite.dat

$(GEOIP_DAT):
	@echo "Downloading geoip.dat..."
	@mkdir -p $(BACKEND_DST)
	curl -L -o $@ $(GEOIP_URL)

$(GEOSITE_DAT):
	@echo "Downloading geosite.dat..."
	@mkdir -p $(BACKEND_DST)
	curl -L -o $@ $(GEOSITE_URL)

.PHONY: build_backend
build_backend: $(NODE_BIN) $(NODEMAN_BIN) $(GEOIP_DAT) $(GEOSITE_DAT)

clean_backend:
	rm -rf $(BACKEND_ROOT)/nodeman/internal/pages/admpage
	rm -rf $(BACKEND_ROOT)/nodeman/internal/pages/userpage
	rm -rf $(BACKEND_DST)

# # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #
# Install
# # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #

INSTALL_PREFIX ?= /usr/local/bin
INSTALL_DATA_DIR ?= /usr/local/share/xray

.PHONY: install
install:
	@echo "Installing binaries to $(INSTALL_PREFIX)..."
	install -m 755 $(BACKEND_DST)/xray-node $(INSTALL_PREFIX)/xray-node
	install -m 755 $(BACKEND_DST)/xray-nodeman $(INSTALL_PREFIX)/xray-nodeman
	@echo "Installing xray data to $(INSTALL_DATA_DIR)..."
	mkdir -p $(INSTALL_DATA_DIR)
	install -m 644 $(GEOIP_DAT) $(INSTALL_DATA_DIR)/geoip.dat
	install -m 644 $(GEOSITE_DAT) $(INSTALL_DATA_DIR)/geosite.dat

.PHONY: uninstall
uninstall:
	@echo "Removing binaries from $(PREFIX)..."
	rm -f $(INSTALL_PREFIX)/xray-node
	rm -f $(INSTALL_PREFIX)/xray-nodeman
	rm -f $(INSTALL_DATA_DIR)/geoip.dat
	rm -f $(INSTALL_DATA_DIR)/geosite.dat