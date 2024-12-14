.PHONY: default
default: help

.PHONY: help
help: ## help information about make commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: run
run: ## start running API
	go run cmd/app/main.go

.PHONY: migrate-new
migrate-new:
	@read -p "Enter the name of the new migration: " name; \
	$(MIGRATE) create -ext sql -dir migrations/ $${name// /_}

.PHONY: migrate-up
migrate-up:
	$(MIGRATE) -path=migrations/ -database "mysql://$(MYSQL_USERNAME):$(MYSQL_PASSWORD)@tcp($(MYSQL_HOST):$(MYSQL_PORT))/$(MYSQL_DATABASE_NAME)" up


.PHONY: migrate-down
migrate-down:
	$(MIGRATE) -path=migrations/ -database "mysql://$(MYSQL_USERNAME):$(MYSQL_PASSWORD)@tcp($(MYSQL_HOST):$(MYSQL_PORT))/$(MYSQL_DATABASE_NAME)" down