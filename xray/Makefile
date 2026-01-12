XRAY_VERSION := 25.6.8
XRAY_URL := https://github.com/XTLS/Xray-core/archive/refs/tags/v$(XRAY_VERSION).zip
GEOIP_URL := https://raw.githubusercontent.com/Loyalsoldier/v2ray-rules-dat/release/geoip.dat
GEOSITE_URL := https://raw.githubusercontent.com/Loyalsoldier/v2ray-rules-dat/release/geosite.dat

BIN_DIR ?= $(CURDIR)/bin
INSTALL_PREFIX ?= /usr/local/bin/xray

# downloading XRay data
XRAY_ZIP := $(CURDIR)/xray.zip
GEOIP_DAT := $(BIN_DIR)/geoip.dat
GEOSITE_DAT := $(BIN_DIR)/geosite.dat

# XRay directory after extraction
XRAY_DIR := $(CURDIR)/Xray-core-$(XRAY_VERSION)

# XRay binary path
XRAY_BIN := $(BIN_DIR)/xray

.PHONY: xray-download xray-extract xray-build xray-clean

xray-download: $(XRAY_ZIP) $(GEOIP_DAT) $(GEOSITE_DAT)

$(XRAY_ZIP):
	@echo "Downloading XRay to $(XRAY_ZIP)..."
	mkdir -p $(CURDIR)
	curl -L -o $(XRAY_ZIP) $(XRAY_URL)
	@echo "XRay downloaded."

$(GEOIP_DAT):
	@echo "Downloading XRay geoip.dat to $(CURDIR)/..."
	curl -L -o $(GEOIP_DAT) $(GEOIP_URL)
	@echo "XRat geoip.dat downloaded."

$(GEOSITE_DAT):
	@echo "Downloading XRay geosite.dat to $(CURDIR)/..."
	curl -L -o $(GEOSITE_DAT) $(GEOSITE_URL)
	@echo "XRat geosite.dat downloaded."

xray-extract: xray-download $(XRAY_DIR)

$(XRAY_DIR):
	@echo "Extracting XRay to $(CURDIR)..."
	unzip -q -o $(XRAY_ZIP) -d $(CURDIR)
	@echo "XRay extracted to $(XRAY_ZIP)."

xray-build: xray-extract $(XRAY_BIN)

$(XRAY_BIN):
	@echo "Build XRay to $(BIN_DIR)..."
	cd $(XRAY_DIR) && GOWORK=off CGO_ENABLED=0 \
		go build -o $(BIN_DIR)/xray \
		-trimpath -buildvcs=false -ldflags="-s -w -buildid=" \
		-v ./main
	@echo "XRay build to $(BIN_DIR)."

xray-install: xray-build
	@echo "Creating XRay install dir: $(INSTALL_PREFIX)"
	mkdir -p $(INSTALL_PREFIX)
	@echo "Installing XRay → $(INSTALL_PREFIX)/..."
	install -m 0755 $(XRAY_BIN) $(INSTALL_PREFIX)/
	@echo "Installing XRay geodata → $(INSTALL_PREFIX)/..."
	install -m 0755 $(GEOIP_DAT) $(INSTALL_PREFIX)/
	install -m 0755 $(GEOSITE_DAT) $(INSTALL_PREFIX)/
	@echo "Xray installed to $(INSTALL_PREFIX)."

xray-clean:
	@echo "Cleaning up XRay build files..."
	rm -rf $(XRAY_DIR) $(XRAY_ZIP) $(GEOIP_DAT) $(GEOSITE_DAT) $(BIN_DIR)/xray
	@echo "XRay build files cleaned up."