.PHONY: terminal web validator test help

env-up: ## Configure the local environment using Docker compose
	docker-compose up

.DEFAULT_GOAL := help

help: ## Prints this help.
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

