ROOT := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
ASSET_OUTPUT = /asset-output
LAMBDA_SRC = backend/pkg/*/*.go backend/pkg/*/*/*.go
LAMBDA_FUNCTIONS = \
	build/apiHandler \
	build/scanRepo \
	build/updateDB

lambda: $(LAMBDA_FUNCTIONS)

build/apiHandler: backend/lambda/apiHandler/*.go $(LAMBDA_SRC)
	go build -o build/apiHandler ./backend/lambda/apiHandler
build/scanRepo: backend/lambda/scanRepo/*.go $(LAMBDA_SRC)
	go build -o build/scanRepo ./backend/lambda/scanRepo
build/updateDB: backend/lambda/updateDB/*.go $(LAMBDA_SRC)
	go build -o build/updateDB ./backend/lambda/updateDB

FRONTEND_DIR = $(ROOT)/frontend
BUNDLE_JS = $(FRONTEND_DIR)/dist/bundle.js
JS_SRC = $(FRONTEND_DIR)/src/js/*.tsx

$(BUNDLE_JS): $(JS_SRC)
	cd $(FRONTEND_DIR) && npm i && npm exec webpack && cd $(ROOT)

js: $(BUNDLE_JS)

asset: lambda js
	cp build/* $(ASSET_OUTPUT)
	cp -r frontend/dist $(ASSET_OUTPUT)/assets
