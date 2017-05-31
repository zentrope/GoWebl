PACKAGE = github.com/zentrope/webl

.PHONY: build-admin build init govendor vendor vendor-check vendor-unused help
.DEFAULT_GOAL := help

# TO update dependencies, govendor fetch +vendor

govendor:
	@hash govendor > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		go get -v -u github.com/kardianos/govendor; \
	fi

ricebox:
	@hash rice > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		go get -v -u github.com/GeertJohan/go.rice/rice; \
	fi

vendor: govendor ## Install govendor and sync deps
	govendor sync

vendor-check: ## Verify that vendored packages match git HEAD
	@git diff-index --quiet HEAD vendor/ || (echo "check-vendor target failed: vendored packages out of sync" && echo && git diff vendor/ && exit 1)

run: vendor ## Run the app from source
	go run main.go

vendor-unused: govendor ## Find unused vendored dependencies
	@govendor list +unused

init: vendor ricebox ## Make sure everything is set up properly for dev.
	cd admin ; yarn

build-admin: ## Build the admin client
	cd admin; yarn ; yarn build

build-freebsd: init build-admin ## Build a version for FreeBSD
	cd internal ; rm -f rice-box.go ;  rice -v embed-go
	GOOS=freebsd GOARCH=amd64 go build -o webl

build: init build-admin ## Build webl into a local binary ./webl.
	cd internal ; rm -f rice-box.go ;  rice -v embed-go
	go build -o webl

clean: ## Clean build artifacts.
	rm -f internal/rice-box.go
	rm -rf webl
	rm -rf admin/build

dist-clean: clean ## Clean everything (vendor, node_modules).
	rm -rf vendor/*/
	rm -rf admin/node_modules

help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-25s\033[0m %s\n", $$1, $$2}' | sort
