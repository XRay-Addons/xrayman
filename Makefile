ROOT := $(CURDIR)
DST := $(ROOT)/build


# # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #
# Build frontend # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #
# # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #

PNPM := pnpm

FRONTEND_ROOT := $(ROOT)/frontend
FRONTEND_DST := $(DST)/frontend

.PHONY: all_f clean_f gen_f build_f

all_f: build_f

gen_f: 
	@echo "Generating frontend..."
	cd $(FRONTEND_ROOT) && $(PNPM) run gen
	cd $(ROOT)

build_f: gen_f
	@echo "Building frontend apps..."
	mkdir -p $(FRONTEND_DST)
	cd $(FRONTEND_ROOT) && $(PNPM) run build
	cp -r $(FRONTEND_ROOT)/admpage/dist $(FRONTEND_DST)/admpage
	cp -r $(FRONTEND_ROOT)/userpage/dist $(FRONTEND_DST)/userpage
	cd $(ROOT)

clean_f:
	rm -rf $(FRONTEND_DST)
	rm -rf $(FRONTEND_ROOT)/admpage/dist
	rm -rf $(FRONTEND_ROOT)/userpage/dist

# # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #
# Build  backend # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #
# # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #

GO := go

BACKEND_ROOT := $(ROOT)/backend
BACKEND_DST := $(DST)/backend

.PHONY: all_b clean_b gen_b embed_frontend_b build_b

all_b: build_b

gen_b: 
	@echo "Generating Backend..."
	cd $(BACKEND_ROOT)/node && go generate ./...
	cd $(BACKEND_ROOT)/nodeman && go generate ./...
	cd $(ROOT)

embed_frontend_b: build_f
	@echo "Embedding Frontend into Backend..."
	mkdir -p $(BACKEND_ROOT)/nodeman/internal/pages
	cp -r $(FRONTEND_DST)/admpage $(BACKEND_ROOT)/nodeman/internal/pages/admpage
	cp -r $(FRONTEND_DST)/userpage $(BACKEND_ROOT)/nodeman/internal/pages/userpage

build_b: gen_b embed_frontend_b
	@echo "Building Backend..."
	mkdir -p $(BACKEND_DST)
	cd $(BACKEND_ROOT) && $(GO) build -o $(BACKEND_DST)/node ./node/cmd/main.go
	cd $(BACKEND_ROOT) && $(GO) build -o $(BACKEND_DST)/nodeman ./nodeman/cmd/main.go
	cd $(ROOT)

clean_b:
	rm -rf $(BACKEND_ROOT)/nodeman/internal/pages/admpage
	rm -rf $(BACKEND_ROOT)/nodeman/internal/pages/userpage
	rm -rf $(BACKEND_DST)