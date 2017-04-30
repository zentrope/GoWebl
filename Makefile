PACKAGE = github.com/zentrope/webl

.PHONY: vendor vendor-check vendor-unused help
.DEFAULT_GOAL := help

vendor: ## Install govendor and sync deps
	go get github.com/kardianos/govendor
	govendor sync ${PACKAGE}

vendor-check: ## Verify that vendored packages match git HEAD
	@git diff-index --quiet HEAD vendor/ || (echo "check-vendor target failed: vendored packages out of sync" && echo && git diff vendor/ && exit 1)

run: vendor ## Run the app in development mode.
	go run main.go

vendor-unused: ## Find unused vendored dependencies
	@hash govendor > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		go get -u github.com/kardianos/govendor; \
	fi
	@govendor list +unused

help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' | sort
