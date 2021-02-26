GO   := go

DIRS_TO_CLEAN:=  ./tmp
FILES_TO_CLEAN:= ./bin/*

ifeq ($(origin GO), undefined)
  GO:=$(shell which go)
endif
ifeq ($(GO),)
  $(error Could not find 'go' in path. Please install go, or if already installed either add it to your path or set GO to point to its directory)
endif

pkgs  = $(shell GOFLAGS=-mod=mod)
pkgDirs = $(shell GOFLAGS=-mod=mod)

GOLANGCI:=$(shell command -v golangci-lint 2> /dev/null)

#-------------------------
# Final targets
#-------------------------
.PHONY: dev

## Build and run
dev.run: dev run

## Execute development pipeline
dev: mod lint swagger build

## Run
run:
	./bin/pokemonb2w

#-------------------------
# Checks
#-------------------------
.PHONY: lint stats.loc

mod: mod.tidy
	$(GO) mod download

mod.tidy:
	$(GO) mod tidy

## Validate code
lint:
ifndef GOLANGCI
	$(error "Please install golangci!")
endif
	@golangci-lint run -v $(pkgDirs)

test:
	$(GO) test -v ./internal/ 

#-------------------------
# Build artefacts
#-------------------------
.PHONY: build build.pokemonb2w

## Build all binaries
build:
	CGO_ENABLED=1 $(GO) build -o  bin/pokemonb2w internal/app.go 

## Compress all binaries
pack:
	@echo ">> packing all binaries"
	@upx -7 -qq bin/*

#-------------------------
# Target: clean
#-------------------------
.PHONY: clean clean.pokemonb2w

## Clean build files
clean:
	rm -rf $(DIRS_TO_CLEAN)
	rm -f $(FILES_TO_CLEAN)

#-------------------------
# Target: swagger
#-------------------------
.PHONY: swagger

swagger: swagger.gen swagger.validate

## Generate swagger json
swagger.gen:
	@echo ">> generating swagger json"
	swagger generate spec -o ./static/swagger.json

## Validate swagger
swagger.validate:
	@echo ">> validating swagger json"
	swagger validate ./static/swagger.json