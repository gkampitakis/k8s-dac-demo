.PHONY: install-tools lint format help docker-image-local deploy cleanup
.SILENT: install-tools lint format help docker-image-local deploy cleanup

help:
	grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

install-tools: ## Install linting tools
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.49.0
	go install mvdan.cc/gofumpt@latest
	go install github.com/segmentio/golines@latest

lint: ## Run golangci linter
	golangci-lint run -c ./golangci.yml ./...

format: ## Format code
	gofumpt -l -w -extra .
	golines . -w

docker-image: ## Build Docker Image for local minikube setup
	docker build . -t localhost:5000/webhook:latest
	docker push localhost:5000/webhook:latest

deploy: ## Deploy all webhook components
	./deploy.sh

cleanup: ## Clean up all webhook components
	./cleanup.sh
