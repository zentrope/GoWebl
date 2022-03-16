##
## Copyright (c) 2017-2022 Keith Irwin
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

DB_PASS = wanheda
DB_USER = webl_user
DB_NAME = webl_db

DB_CREATE = create database $(DB_NAME) with encoding 'UTF8'
DB_SETUP = create user $(DB_USER) with login password '$(DB_PASS)' ;\
	alter database $(DB_NAME) owner to $(DB_USER) ;\
	create extension if not exists pgcrypto

.DEFAULT_GOAL := help

.PHONY: psqldep

PSQL = psql

psqldep:
	@hash $(PSQL) > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		echo "$(PSQL) not found. Try 'brew install postgresql'."; \
		exit 1; \
	fi

.PHONY: db-clean db-init

db-clean: psqldep ## Delete the local dev database
	$(PSQL) template1 -c "drop database $(DB_NAME)"
	$(PSQL) template1 -c "drop user $(DB_USER)"

db-init: psqldep ## Create a local dev database with default creds
	$(PSQL) template1 -c "$(DB_CREATE)"
	$(PSQL) $(DB_NAME) -c "$(DB_SETUP)"

db-schema: psqldep ## Try to use the webl server to load the schema
	@cd ../WeblServer ; go run . -resources resources -assets assets & echo $$! > /tmp/webl.pid
	@sleep 2
	@kill $$(cat /tmp/webl.pid)
	@cd ../WeblData

.PHONY: help

help: ## Display this help listing
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-25s\033[0m %s\n", $$1, $$2}' \
		| sort
