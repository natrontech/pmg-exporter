.PHONY: help run-portal run-orchestrator

help: ## Display this help message
	@cat $(MAKEFILE_LIST) | grep -e "^[a-zA-Z_-]*: *.*## *" | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

run-portal: ## Run the portal
	cd services/portal && yarn dev

run-orchestrator: ## Run the orchestrator
	@cd services/orchestrator && \
	set -a && source .env && set +a && \
	air

generate-graphql: ## Generate GraphQL code for the orchestrator service
	@cd services/orchestrator && \
	go generate ./graph/resolver.go
