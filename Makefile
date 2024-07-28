.PHONY: help lint go-update build docker down run clean

help: ## Display this help message
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

lint: ## Run linting
	./scripts/lint_go.sh

go-update: ## Update all Go dependencies
	./scripts/update_all.sh

build: ## Build the project
	go build -o /bot ./cmd/main.go ./cmd/bot.go

docker: ## Build and run the project in a Docker container in bg
	sudo docker compose up --build -d

down: ## Delete audio-files and stop services
	rm -rf /app/audio-files
	docker compose down --remove-orphans
run: ## Run the project locally
	/bot

clean: ## Clean up the project
	rm -rf /bot