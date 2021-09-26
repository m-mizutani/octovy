ROOT_DIR := $(dir $(abspath $(firstword $(MAKEFILE_LIST))))

ASSET_DIR=./assets
DIST_DIR=$(ASSET_DIR)/dist
ASSET_OUT=$(DIST_DIR)/out/index.html
ASSET_SRC=$(ASSET_DIR)/**/*.tsx

SRC=pkg/**/*.go
ENT_DIR=./pkg/infra/ent
ENT_SRC=$(ENT_DIR)/ent.go
ENT_SCHEMA_DIR=./pkg/domain/schema

BINARY=./octovy
EXAMPLE_SRC_DIR=./examples/basic

IGNORE_TEST=\
	"" \
	"^github.com/m-mizutani/octovy/pkg/infra/domain/types" \
	"^github.com/m-mizutani/octovy/pkg/infra" \

all: $(BINARY)

ent: $(ENT_SRC)

asset: $(ASSET_OUT)

$(ASSET_OUT): $(ASSET_SRC)
	cd $(ASSET_DIR) && npm run export && cd $(ROOT_DIR)

$(ENT_SRC): $(ENT_SCHEMA_DIR)/*.go
	ent generate $(ENT_SCHEMA_DIR) --target $(ENT_DIR) --feature sql/upsert

dev: $(SRC)
	go run ./cmd/octovy/ serve -d "root:${MYSQL_ROOT_PASSWORD}@tcp(localhost:3306)/${MYSQL_DATABASE}"

test: $(SRC) $(ENT_SRC)
	go list ./... | \
	grep -v "^github.com/m-mizutani/octovy/pkg/infra/ent$$" | \
	grep -v "^github.com/m-mizutani/octovy/pkg/infra/ent/" | \
	grep -v "^github.com/m-mizutani/octovy/pkg/infra/domain/types$$" | \
	grep -v "^github.com/m-mizutani/octovy/pkg/infra$$" | \
	xargs go test

octovy: $(SRC) $(ENT_SRC) $(ASSET_OUT)
	go build -o $(BINARY) ./cmd/octovy

clean:
	rm -f $(ENT_SRC)
	rm -f $(BINARY)
