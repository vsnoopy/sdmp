.PHONY: help

help: ## Display this help message
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

lint: ## Run linting
	./scripts/lint_go.sh

go-update: ## Update all Go dependencies
	./scripts/update_all.sh

build: ## Build the project
	go build -o /bot ./cmd/main.go ./cmd/bot.go

run: ## Run the project locally
	/bot

clean: ## Clean up the project
	rm -rf /bot