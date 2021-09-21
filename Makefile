ROOT_DIR := $(dir $(abspath $(firstword $(MAKEFILE_LIST))))

ASSET_DIR=./assets
DIST_DIR=$(ASSET_DIR)/dist
ASSET_JS=$(DIST_DIR)/bundle.js
ASSET_SRC=$(ASSET_DIR)/**/*.tsx
ASSETS=$(ASSET_JS) $(DIST_DIR)/index.html

SRC= \
	cmd/octovy/*.go \
	pkg/**/*.go
ENT_DIR=./pkg/infra/ent
ENT_SRC=$(ENT_DIR)/ent.go
ENT_SCHEMA_DIR=./pkg/domain/schema

BINARY=./octovy
EXAMPLE_SRC_DIR=./examples/basic

all: $(BINARY)

ent: $(ENT_SRC)

docker:
	docker run -p 127.0.0.1:3306:3306 -e MYSQL_ROOT_PASSWORD -e MYSQL_DATABASE mysql

$(ASSET_JS): $(ASSET_SRC)
	cd $(ASSET_DIR) && npm run export && cd $(ROOT_DIR)

$(ENT_SRC): $(ENT_SCHEMA_DIR)/*.go
	ent generate $(ENT_SCHEMA_DIR) --target $(ENT_DIR) --feature sql/upsert

dev: $(SRC)
	go run ./cmd/octovy/ serve -d "root:${MYSQL_ROOT_PASSWORD}@tcp(localhost:3306)/${MYSQL_DATABASE}"

test: $(SRC) $(ENT_SRC)
	go test ./...

$(CHAIN): $(EXAMPLE_SRC_DIR)/*.go $(SRC) $(ENT_SRC)
	go build -buildmode=plugin -o chain.so $(EXAMPLE_SRC_DIR)

octovy: $(SRC) $(ENT_SRC)
	go build -o $(BINARY) ./cmd/octovy

clean:
	rm -f $(ENT_SRC)
	rm -f $(BINARY)
