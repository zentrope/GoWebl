##
## Copyright (c) 2017 Keith Irwin
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

PACKAGE = github.com/zentrope/webl

DB_PASS = wanheda
DB_USER = webl_user
DB_NAME = webl_db

DB_CREATE = create database $(DB_NAME) with encoding 'UTF8'
DB_SETUP = create user $(DB_USER) with login password '$(DB_PASS)' ;\
	alter database $(DB_NAME) owner to $(DB_USER) ;\
	create extension if not exists pgcrypto

.PHONY: build-admin build init godep vendor vendor-check vendor-unused help
.PHONY: db-clean db-init

.DEFAULT_GOAL := help


##-----------------------------------------------------------------------------
## Make depenencies
##-----------------------------------------------------------------------------

godep:
	@hash dep > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		go get -v -u github.com/golang/dep/cmd/dep; \
	fi

##-----------------------------------------------------------------------------
## Project dependencies
##-----------------------------------------------------------------------------

vendor: godep ## Install and sync deps
	dep ensure

init: vendor ## Make sure everything is set up properly for dev.
	cd admin ; yarn

##-----------------------------------------------------------------------------

run: vendor ## Run the app from source
	go run main.go

build-admin: ## Build the admin client
	cd admin; yarn ; yarn build

build-freebsd: init build-admin ## Build a version for FreeBSD
	GOOS=freebsd GOARCH=amd64 go build -o webl

build: init build-admin ## Build webl into a local binary ./webl.
	go build -o webl

clean: ## Clean build artifacts.
	rm -rf webl
	rm -rf admin/build

dist-clean: clean ## Clean everything (vendor, node_modules).
	rm -rf vendor
	rm -rf admin/node_modules

db-clean: ## Delete the local dev database
	psql template1 -c "drop database $(DB_NAME)"
	psql template1 -c "drop user $(DB_USER)"

db-init: ## Create a local dev database with default creds
	psql template1 -c "$(DB_CREATE)"
	psql $(DB_NAME) -c "$(DB_SETUP)"

help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-25s\033[0m %s\n", $$1, $$2}' \
		| sort
