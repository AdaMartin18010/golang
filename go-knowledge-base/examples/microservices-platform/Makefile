.PHONY: help build up down test lint clean

help:
	@echo "Microservices Platform - Available Commands:"
	@echo "  build        - Build all services"
	@echo "  up           - Start all services with Docker Compose"
	@echo "  down         - Stop all services"
	@echo "  test         - Run all tests"
	@echo "  lint         - Run linters"
	@echo "  clean        - Clean up Docker resources"

build:
	@echo "Building all services..."
	cd services/user-service && docker build -t user-service .
	cd services/order-service && docker build -t order-service .
	cd services/payment-service && docker build -t payment-service .
	cd services/inventory-service && docker build -t inventory-service .
	cd services/notification-service && docker build -t notification-service .

up:
	cd deployments && docker-compose up -d

down:
	cd deployments && docker-compose down

test:
	@echo "Running tests..."
	cd services/user-service && go test ./...
	cd services/order-service && go test ./...

lint:
	@echo "Running linters..."
	golangci-lint run ./...

clean:
	cd deployments && docker-compose down -v --rmi local

docker-up:
	cd deployments && docker-compose up -d --build

k8s-deploy:
	kubectl apply -f infra/k8s/

k8s-delete:
	kubectl delete -f infra/k8s/ --ignore-not-found
