.PHONY: up down build rebuild generate help
.DEFAULT_GOAL := help

up: ## Start the docker containers
	docker-compose up

down: ## Stop the docker containers
	docker-compose down

build: ## Build the docker containers
	docker-compose build

rebuild: down generate build up ## Rebuild templates and containers

generate: ## Generate templ templates
	templ generate

help: ## Display this help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)