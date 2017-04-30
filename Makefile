PACKAGE = github.com/zentrope/webl

.PHONY: vendor check-vendor help
.DEFAULT_GOAL := help

vendor: ## Install govendor and sync deps
	go get github.com/kardianos/govendor
	govendor sync ${PACKAGE}

check-vendor: ## Verify that vendored packages match git HEAD
	@git diff-index --quiet HEAD vendor/ || (echo "check-vendor target failed: vendored packages out of sync" && echo && git diff vendor/ && exit 1)

help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
