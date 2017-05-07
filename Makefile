PACKAGE = github.com/zentrope/webl

.PHONY: admin build govendor vendor vendor-check vendor-unused help
.DEFAULT_GOAL := help

govendor:
	@hash govendor > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		go get -u github.com/kardianos/govendor; \
	fi

vendor: govendor ## Install govendor and sync deps
	govendor sync

admin: ## Build the admin client
	cd admin; yarn ; yarn build

vendor-check: ## Verify that vendored packages match git HEAD
	@git diff-index --quiet HEAD vendor/ || (echo "check-vendor target failed: vendored packages out of sync" && echo && git diff vendor/ && exit 1)

run: vendor ## Run the app from source
	go run main.go

vendor-unused: govendor ## Find unused vendored dependencies
	@govendor list +unused

init: vendor ## Make sure everything is set up properly for dev.
	cd admin ; yarn

build: vendor admin ## Build webl into a local binary ./webl.
	go build

clean: ## Clean build artifacts.
	rm -rf webl
	rm -rf admin/build

dist-clean: clean ## Clean everything (vendor, node_modules).
	rm -rf vendor/*/
	rm -rf admin/node_modules

help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-25s\033[0m %s\n", $$1, $$2}' | sort
