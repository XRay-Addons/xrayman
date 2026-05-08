SHELL := /bin/bash

BIN_DIR := $(CURDIR)/bin
NPM_ROOT := $(CURDIR)/nodeman/web

USERPAGE_WEB_SRC := $(CURDIR)/nodeman/web/pages/userpage
ADMPAGE_WEB_SRC := $(CURDIR)/nodeman/web/pages/admpage

USERPAGE_WEB_DST := $(CURDIR)/nodeman/internal/pages/userpage
ADMPAGE_WEB_DST := $(CURDIR)/nodeman/internal/pages/admpage


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


.PHONY: build
build: npm-install userpage admpage
	@echo "Frontend build done."