.PHONY: help build up down logs test lint clean k8s-deploy k8s-delete

# Variables
DOCKER_COMPOSE = docker-compose
KUBECTL = kubectl
NAMESPACE = microservices

help: ## Show this help
	@echo "Microservices Platform - Available Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build all services
	@echo "Building all services..."
	$(DOCKER_COMPOSE) build --parallel

up: ## Start all services
	@echo "Starting all services..."
	$(DOCKER_COMPOSE) up -d

down: ## Stop all services
	@echo "Stopping all services..."
	$(DOCKER_COMPOSE) down

logs: ## View logs for all services
	$(DOCKER_COMPOSE) logs -f

logs-%: ## View logs for a specific service (e.g., make logs-user-service)
	$(DOCKER_COMPOSE) logs -f $*

ps: ## List running services
	$(DOCKER_COMPOSE) ps

test: ## Run tests for all services
	@echo "Running tests..."
	cd ../services/user-service && go test ./...
	cd ../services/order-service && go test ./...

test-integration: ## Run integration tests
	@echo "Running integration tests..."
	$(DOCKER_COMPOSE) -f docker-compose.yml -f docker-compose.test.yml up --abort-on-container-exit

lint: ## Run linter on all services
	@echo "Running linters..."
	cd ../services/user-service && golangci-lint run
	cd ../services/order-service && golangci-lint run

clean: ## Clean up Docker resources
	@echo "Cleaning up..."
	$(DOCKER_COMPOSE) down -v --rmi local
	docker system prune -f

# Kubernetes commands
k8s-deploy: ## Deploy to Kubernetes
	@echo "Deploying to Kubernetes..."
	$(KUBECTL) create namespace $(NAMESPACE) --dry-run=client -o yaml | $(KUBECTL) apply -f -
	$(KUBECTL) apply -f ../infra/k8s/ -n $(NAMESPACE)

k8s-delete: ## Delete from Kubernetes
	@echo "Deleting from Kubernetes..."
	$(KUBECTL) delete -f ../infra/k8s/ -n $(NAMESPACE) --ignore-not-found=true

k8s-status: ## Check Kubernetes status
	$(KUBECTL) get pods -n $(NAMESPACE)
	$(KUBECTL) get svc -n $(NAMESPACE)

k8s-logs: ## View Kubernetes logs
	$(KUBECTL) logs -l app=user-service -n $(NAMESPACE) --tail=100 -f

k8s-port-forward: ## Port forward services
	@echo "Port forwarding services..."
	$(KUBECTL) port-forward svc/user-service 8081:80 -n $(NAMESPACE) &
	$(KUBECTL) port-forward svc/order-service 8082:80 -n $(NAMESPACE) &

# Service-specific commands
user-service-build: ## Build user service
	docker build -t microservices/user-service:latest ../services/user-service

user-service-run: ## Run user service locally
	cd ../services/user-service && go run cmd/main.go

order-service-build: ## Build order service
	docker build -t microservices/order-service:latest ../services/order-service

order-service-run: ## Run order service locally
	cd ../services/order-service && go run cmd/main.go

# Load testing
load-test: ## Run load tests
	k6 run ../tests/load/order_creation.js

# Monitoring
prometheus: ## Open Prometheus
	open http://localhost:9090

grafana: ## Open Grafana
	open http://localhost:3000

jaeger: ## Open Jaeger
	open http://localhost:16686

kafka-ui: ## Open Kafka UI
	open http://localhost:8090
