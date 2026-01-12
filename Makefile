BIN_DIR := $(CURDIR)/bin
INSTALL_PREFIX := /usr/local/bin/xrayman

NODE_SRC := ./node/cmd
NODEMAN_SRC := ./nodeman/cmd

NODE_BIN := $(BIN_DIR)/node
NODEMAN_BIN := $(BIN_DIR)/nodeman
XRAY_BIN := $(BIN_DIR)/xray

.PHONY: all build install clean

all: build

build: $(NODE_BIN) $(NODEMAN_BIN) $(XRAY_BIN) 

$(NODE_BIN):
	@echo "Building node..."
	mkdir -p $(BIN_DIR)
	go build -o $(NODE_BIN) $(NODE_SRC)
	@echo "Node built."

$(NODEMAN_BIN):
	@echo "Building nodeman..."
	mkdir -p $(BIN_DIR)
	go build -o $(NODEMAN_BIN) $(NODEMAN_SRC)
	@echo "Nodeman built."

$(XRAY_BIN):
	@echo "Building XRay to $(BIN_DIR)..."
	mkdir -p $(BIN_DIR)
	make -C xray BIN_DIR=$(BIN_DIR) xray-build 
	@echo "XRay built."

install: build
	@echo "Creating install dir: $(INSTALL_PREFIX)"
	mkdir -p $(INSTALL_PREFIX)

	@echo "Installing node → $(INSTALL_PREFIX)/"
	install -m 0755 $(NODE_BIN) $(INSTALL_PREFIX)/

	@echo "Installing nodeman → $(INSTALL_PREFIX)/"
	install -m 0755 $(NODEMAN_BIN) $(INSTALL_PREFIX)/

	@echo "Installing xray → $(INSTALL_PREFIX)/"
	make -C xray \
		BIN_DIR=$(BIN_DIR) \
		CUR_DIR=$(CURDIR)/xray \
		INSTALL_PREFIX=$(INSTALL_PREFIX) \
		xray-install

	@echo "Installed to $(INSTALL_PREFIX)"

clean:
	@echo "Cleaning up build files..."
	rm -rf $(BIN_DIR)

	make -C xray \
		BIN_DIR=$(BIN_DIR) \
		CUR_DIR=$(CURDIR)/xray \
		INSTALL_PREFIX=$(INSTALL_PREFIX) \
		xray-clean

	@echo "Build files cleaned."
