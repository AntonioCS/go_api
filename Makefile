.DEFAULT_GOAL := help
.PHONY: help

dc-file := ./docker/docker-compose.yaml
env := project/.env

project_dir := project/
project_main := $(project_dir)cmd/app/main.go
project_main_in_container = cmd/app/main.go

project_name := go_api

dc := @docker-compose --file $(dc-file) --env-file $(env) -p $(project_name)
service_go = go_api_app
service_db = go_api_mysql

start: ## Start containers
	$(dc) up -d $(service_db)

stop: ## Stop containers
	$(dc) stop $(service_db)

restart: stop
restart: start
restart: ## Restart all containers
	@echo Finished restart

start-no-d: ## Start containers not in detach mode
	$(dc) up

logs: ## Show container logs
	$(dc) logs -f -t

status: ## List containers
	$(dc) ps

rebuild: ## Rebuild containers
	$(dc) build --pull

db-bash: ## db bash
	@$(dc) exec $(service_db) sh

db-log-status: ## Show logging status
	@$(dc) exec $(service_db) mysql -u $(db_user) -p$(db_pass) -e "SHOW GLOBAL VARIABLES LIKE 'general_log';"

go-bash: ## Bash into Go container
	$(dc) run $(service_go) sh || true

go-run: ## Run program
	$(dc) run --service-ports $(service_go) go run $(project_main_in_container) || true

go-test-api: ## Run the API tests
	$(dc) run -e CGO_ENABLED=0 $(service_go) go test ./cmd/app/ || true

#go-build: ## Build program
#	$(dc) run --service-ports $(service_go) go build -o app $(project_main_in_container) || true

go-tests: ## Run tests
	$(dc) run $(service_go) go test ./pkg/GuestBook/

go-deploy: ## Run to create image
	docker build -f ./docker/deploy/Dockerfile .

.SILENT:
help:
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-40s\033[0m %s\n", $$1, $$2}'
