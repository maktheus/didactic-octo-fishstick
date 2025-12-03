SHELL := /bin/sh

APP         ?= back-end-tcc
PKG         ?= ./cmd/api
BIN_DIR     ?= bin
HTTP_PORT   ?= 8080
CONTAINER   ?= back-end-tcc-dev

.PHONY: help run dev build test docker-build docker-run docker-stop swagger

help: ## List available targets
	@printf "Available targets:\n"
	@grep -E '^[a-zA-Z_-]+:.*##' $(lastword $(MAKEFILE_LIST)) | awk -F':|##' '{printf "  make %-15s %s\n", $$1, $$3}'

run: ## Run the API gateway with current environment variables
	@go run $(PKG)

dev: ## Run the API gateway with development defaults
	@APP_ENV=development HTTP_PORT=$(HTTP_PORT) go run $(PKG)

build: ## Compile the API binary into ./bin/api
	@mkdir -p $(BIN_DIR)
	@CGO_ENABLED=0 go build -trimpath -o $(BIN_DIR)/api $(PKG)

test: ## Execute all Go tests
	@go test ./...

docker-build: ## Build the Docker image (make docker-build APP=name tag)
	@docker build -t $(APP) .

docker-run: ## Run the Docker image exposing HTTP_PORT locally
	@docker run --rm -d --name $(CONTAINER) -p $(HTTP_PORT):8080 $(APP)

docker-stop: ## Stop the running container started with docker-run
	-@docker stop $(CONTAINER) >/dev/null 2>&1 || true

swagger: ## Launch Swagger UI serving docs/openapi.json on http://localhost:8081
	@docker run --rm -p 8081:8080 \
		-e SWAGGER_JSON=/docs/swagger.json \
		-v "$(PWD)"/docs/openapi.json:/docs/swagger.json \
		swaggerapi/swagger-ui
