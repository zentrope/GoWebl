##
## Copyright (c) 2017-present Keith Irwin
##
## This program is free software: you can redistribute it and/or modify
## it under the terms of the GNU General Public License as published
## by the Free Software Foundation, either version 3 of the License,
## or (at your option) any later version.
##
## This program is distributed in the hope that it will be useful,
## but WITHOUT ANY WARRANTY; without even the implied warranty of
## MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
## GNU General Public License for more details.
##
## You should have received a copy of the GNU General Public License
## along with this program.  If not, see <http://www.gnu.org/licenses/>.

.DEFAULT_GOAL := help

##-----------------------------------------------------------------------------
## Make dependencies
##-----------------------------------------------------------------------------

.PHONY: treedep

TREE = tree

treedep:
	@hash $(TREE) > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		echo "$(TREE) not found. Try 'brew install $(TREE)'."; \
		exit 1; \
	fi

##-----------------------------------------------------------------------------
## Project dependencies
##-----------------------------------------------------------------------------

.PHONY: vendor init

vendor: ## Install and sync deps
	go mod download
	go mod tidy

init: ## Make sure everything is set up properly for dev.
	@$(MAKE) vendor

##-----------------------------------------------------------------------------
## Build
##-----------------------------------------------------------------------------

.PHONY: build-freebsd build clean

clean: ## Clean build artifacts.
	rm -rf webl
	rm -rf dist

build-freebsd: init ## Build a version for FreeBSD
	CGO_ENABLED=0 GOOS=freebsd GOARCH=amd64 go build -o webl

build: init ## Build webl into a local binary ./webl.
	CGO_ENABLED=0 go build -o webl

##-----------------------------------------------------------------------------
## Distribute
##-----------------------------------------------------------------------------

.PHONY: dist-prepare dist dist-assemble dist-freebsd dist-clean

DIST = ./dist
DIST_RESOURCES = $(DIST)/resources
DIST_ASSETS = $(DIST)/assets

dist-clean: clean ## Clean everything (vendor, node_modules, dist).
	rm -rf vendor

dist-prepare:
	if [ -e "dist" ]; then rm -rf dist ; fi
	mkdir -p $(DIST_RESOURCES)
	mkdir -p $(DIST_ASSETS)

dist-assemble:
	cp -r resources/* $(DIST_RESOURCES)
	cp -r assets/* $(DIST_ASSETS)
	cp -r webl $(DIST)

dist: ## Build distribution for current platform.
	@$(MAKE) dist-prepare
	@$(MAKE) build
	@$(MAKE) dist-assemble

dist-freebsd: ## Build distribution for FreeBSD.
	@$(MAKE) dist-prepare
	@$(MAKE) build-freebsd
	@$(MAKE) dist-assemble

##-----------------------------------------------------------------------------
## Utilties
##-----------------------------------------------------------------------------

.PHONY: run tree help

run: vendor ## Run the app from source
	go run main.go -resources resources -assets assets || true

tree: treedep ## View source hierarchy without vendor pkgs
	$(TREE) -C -I "node_modules|vendor|build|dist" || true

help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-25s\033[0m %s\n", $$1, $$2}' \
		| sort
