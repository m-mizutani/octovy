ROOT := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

CMD=octovy
BACKEND_DIR=pkg
FRONTEND_DIR=assets
ASSET=$(FRONTEND_DIR)/dist/bundle.js

all: $(CMD)

$(ASSET): $(FRONTEND_DIR)/src/**/*.tsx $(FRONTEND_DIR)/src/*.tsx
	cd $(ROOT)/$(FRONTEND_DIR) && npm install && cd $(ROOT)

$(CMD): $(ASSET) $(BACKEND_DIR)/**/*.go
	go build -v -o $(CMD) ./cmd/octovy
