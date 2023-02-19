include .env

.PHONY: help
help: ## Display this help screen. This Makefile is prepared to build the environment for operation confirmation on the local machine.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PNONY: env
env: ## Create .env file.
	@cp .env.local .env

.PHONY: up
up: ## Create development environment.
	@docker compose --project-name ${APP_NAME} up -d

.PHONY: down
down: ## Delete development environment.
	@docker compose --project-name ${APP_NAME} down

.PHONY: psql
psql: ## Connect local postgres.
	@docker exec -it ${APP_NAME}-postgres psql -U postgres
