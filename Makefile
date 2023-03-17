include .env
export

.PHONY: help
help: ## display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

# develop

.PNONY: env
env: ## Create .env file.
	@cp .env.local .env

.PHONY: aqua
aqua: ## Put the path in your environment variables. ex) export PATH="${AQUA_ROOT_DIR:-${XDG_DATA_HOME:-$HOME/.local/share}/aquaproj-aqua}/bin:$PATH"
	@go run github.com/aquaproj/aqua-installer@latest --aqua-version v1.33.0

.PHONY: tool
tool: ## Install tool.
	@aqua i
	@(cd api && npm install)

.PHONY: dev
dev: ## Make development.
	@docker compose --project-name ${APP_NAME} --file ./.docker/docker-compose.yaml up -d

.PHONY: redev
redev: ## restart dev container
	@touch cmd/app/core/main.go
	@touch cmd/app/api/main.go

.PHONY: down
down: ## Down development. (retain containers and delete volumes.)
	@docker compose --project-name ${APP_NAME} down --volumes

.PHONY: balus
balus: ## Destroy everything about docker. (containers, images, volumes, networks.)
	@docker compose --project-name ${APP_NAME} down --rmi all --volumes

.PHONY: psql
psql: ## Connect to postgres.
	@docker exec -it ${APP_NAME}-postgres psql -U postgres

.PHONY: migrate
migrate: ## Migrate database.
	@touch cmd/db/migrate/main.go

.PHONY: backup
backup: ## Backup database.
	@touch cmd/db/backup/main.go

# go

.PHONY: fmt
fmt: ## Format code.
	@go fmt ./...

.PHONY: lint
lint: ## Lint code.
	@golangci-lint run --fix

.PHONY: mod
mod: ## Go mod
	@go mod tidy
	@go mod vendor

.PHONY: modules
modules: ## List modules with dependencies.
	@go list -u -m all

.PHONY: renovate
renovate: ## Update modules with dependencies.
	@go get -u -t ./...
	@go mod tidy
	@go mod vendor

.PHONY: compile
compile: ## Compile code.
	@go build -v ./... && go clean

.PHONY: test
test: ## Run unit test. If you want to invalidate the cache, please specify an argument like `make test c=c`.
	@$(call _test,${c})

define _test
if [ -z "$1" ]; then \
	go test ./internal/... ; \
else \
	go test ./internal/... -count=1 ; \
fi
endef

.PHONY: e2e
e2e: ## Run e2e test. If you want to invalidate the cache, please specify an argument like `make e2e c=c`.
	@$(call _e2e,${c})

define _e2e
if [ -z "$1" ]; then \
	go test ./e2e/... ; \
else \
	go test ./e2e/... -count=1 ; \
fi
endef

.PHONY: gen
gen: ## Generate code.
	@go generate ./...
	@oapi-codegen -generate types -package openapi ./api/openapi.yaml > ./pkg/openapi/types.gen.go
	@oapi-codegen -generate chi-server -package openapi ./api/openapi.yaml > ./pkg/openapi/server.gen.go
	@oapi-codegen -generate client -package openapi ./api/openapi.yaml > ./pkg/openapi/client.gen.go
	@(cd proto && buf generate --template buf.gen.yaml)
	@go mod tidy

# support

.PHONY: doc
doc: ## Generate documentation.
	@rm -rf doc
	@mkdir -p doc/proto
	@mkdir -p doc/openapi
	@tbls doc $(DATABASE_URL) doc/databases
	@protoc --doc_out=./doc/proto --doc_opt=markdown,README.md proto/article/**/*.proto proto/health/**/*.proto
	@(cd api && npx widdershins --omitHeader --code true openapi.yaml ../doc/openapi/README.md)

.PHONY: buflint
buflint: ## Lint proto file.
	@(cd proto && buf lint)

.PHONY: bufmt
bufmt: ## Format proto file.
	@(cd proto && buf format -w)

.PHONY: apilint
apilint: ## Lint api file.
	@(cd api && npx spectral lint openapi.yaml)

.PHONY: ymlint
ymlint: ## Lint yaml file.
	@yamlfmt -lint

.PHONY: ymlfmt
ymlfmt: ## Format yaml file.
	@yamlfmt
