# Microservices Platform Example

A comprehensive, production-ready microservices platform built with Go, demonstrating modern cloud-native architecture patterns, service mesh integration, and enterprise-grade observability.

## Table of Contents

1. [Overview](#overview)
2. [Architecture](#architecture)
3. [Services](#services)
4. [Getting Started](#getting-started)
5. [Development Guide](#development-guide)
6. [Deployment](#deployment)
7. [Observability](#observability)
8. [Security](#security)
9. [Performance](#performance)
10. [Troubleshooting](#troubleshooting)

## Overview

This microservices platform demonstrates a complete e-commerce system with the following characteristics:

- **10+ Microservices**: User Service, Order Service, Payment Service, Inventory Service, Notification Service, API Gateway, and more
- **Service Mesh**: Istio/Linkerd integration for traffic management
- **Event-Driven**: Async communication via Kafka/NATS
- **Polyglot Persistence**: PostgreSQL, MongoDB, Redis, Elasticsearch
- **Observability**: Distributed tracing, metrics, structured logging
- **Security**: mTLS, JWT authentication, RBAC, secrets management
- **Resilience**: Circuit breakers, retries, timeouts, bulkheads
- **Scalability**: Horizontal pod autoscaling, load balancing

### Technology Stack

| Component | Technology |
|-----------|------------|
| Language | Go 1.21+ |
| Framework | Gin/Echo/Fiber |
| RPC | gRPC with Protocol Buffers |
| Message Broker | Apache Kafka, NATS |
| Service Mesh | Istio |
| API Gateway | Kong/Envoy |
| Databases | PostgreSQL 15, MongoDB 6, Redis 7 |
| Search | Elasticsearch 8 |
| Monitoring | Prometheus, Grafana, Jaeger |
| Container | Docker, Kubernetes |
| CI/CD | GitHub Actions, ArgoCD |

## Architecture

### High-Level Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              Client Layer                                    │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐ │
│  │   Web App   │  │ Mobile App  │  │  Admin UI   │  │   Third-party API   │ │
│  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘  └──────────┬──────────┘ │
└─────────┼────────────────┼────────────────┼────────────────────┼────────────┘
          │                │                │                    │
          └────────────────┴────────────────┴────────────────────┘
                                   │
                    ┌──────────────▼──────────────┐
                    │      CDN / WAF / DDoS       │
                    │         CloudFlare          │
                    └──────────────┬──────────────┘
                                   │
┌──────────────────────────────────▼──────────────────────────────────────────┐
│                           API Gateway Layer                                  │
│  ┌────────────────────────────────────────────────────────────────────────┐  │
│  │                         Kong / Envoy Gateway                            │  │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐   │  │
│  │  │   Rate      │  │   Auth      │  │   Routing   │  │  Transform  │   │  │
│  │  │  Limiting   │  │   JWT/OAuth │  │   Rules     │  │  Plugins    │   │  │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘   │  │
│  └────────────────────────────────────────────────────────────────────────┘  │
└──────────────────────────────────┬──────────────────────────────────────────┘
                                   │
┌──────────────────────────────────▼──────────────────────────────────────────┐
│                          Service Mesh (Istio)                                │
│  ┌────────────────────────────────────────────────────────────────────────┐  │
│  │                    Virtual Services & Destination Rules                 │  │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐   │  │
│  │  │   mTLS      │  │  Traffic    │  │   Circuit   │  │   Retry     │   │  │
│  │  │   Policy    │  │   Split     │  │   Breaker   │  │   Policy    │   │  │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘   │  │
│  └────────────────────────────────────────────────────────────────────────┘  │
└──────────────────────────────────┬──────────────────────────────────────────┘
                                   │
          ┌────────────────────────┼────────────────────────┐
          │                        │                        │
┌─────────▼──────────┐  ┌──────────▼──────────┐  ┌──────────▼──────────┐
│   Core Services    │  │   Support Services  │  │   Event Bus         │
│                    │  │                     │  │                     │
│ ┌───────────────┐  │  │ ┌───────────────┐   │  │ ┌───────────────┐   │
│ │  User Service │  │  │ │ Config Server│   │  │ │     Kafka     │   │
│ │  (Go/Gin)     │  │  │ │  (Go/Etcd)   │   │  │ │  (3 brokers)  │   │
│ └───────────────┘  │  │ └───────────────┘   │  │ └───────────────┘   │
│ ┌───────────────┐  │  │ ┌───────────────┐   │  │ ┌───────────────┐   │
│ │ Order Service │  │  │ │  Discovery   │   │  │ │     NATS      │   │
│ │  (Go/Echo)    │  │  │ │  (Consul)    │   │  │ │  (Streaming)  │   │
│ └───────────────┘  │  │ └───────────────┘   │  │ └───────────────┘   │
│ ┌───────────────┐  │  │ ┌───────────────┐   │  └─────────────────────┘
│ │Payment Service│  │  │ │  API Gateway │   │
│ │  (Go/gRPC)    │  │  │ │   (Kong)     │   │
│ └───────────────┘  │  │ └───────────────┘   │
│ ┌───────────────┐  │  │ ┌───────────────┐   │
│ │Inventory Svc  │  │  │ │   Gateway    │   │
│ │  (Go/Fiber)   │  │  │ │  (Envoy)     │   │
│ └───────────────┘  │  │ └───────────────┘   │
│ ┌───────────────┐  │  └─────────────────────┘
│ │Notification Svc│  │
│ │  (Go/NATS)    │  │
│ └───────────────┘  │
│ ┌───────────────┐  │
│ │Shipping Service│  │
│ │  (Go/gRPC)    │  │
│ └───────────────┘  │
│ ┌───────────────┐  │
│ │Analytics Svc  │  │
│ │  (Go/Kafka)   │  │
│ └───────────────┘  │
└─────────┬──────────┘
          │
    ┌─────┴─────┬─────────────┬─────────────┐
    │           │             │             │
┌───▼───┐  ┌───▼────┐   ┌────▼────┐  ┌────▼────┐
│PostgreSQL│  │MongoDB │   │  Redis  │  │Elasticsearch│
│ (Users)  │  │(Orders)│   │ (Cache) │  │  (Search)   │
└─────────┘  └────────┘   └─────────┘  └─────────────┘
```

### Service Interaction Flow

```
┌─────────┐     ┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│  Client │────▶│ API Gateway │────▶│  User Svc   │────▶│ PostgreSQL  │
└─────────┘     └─────────────┘     └─────────────┘     └─────────────┘
                                           │
                                           ▼
                                    ┌─────────────┐
                                    │    Redis    │
                                    │   (Cache)   │
                                    └─────────────┘

Order Creation Flow:
┌─────────┐     ┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│  Client │────▶│ API Gateway │────▶│  Order Svc  │────▶│  MongoDB    │
└─────────┘     └─────────────┘     └──────┬──────┘     └─────────────┘
                                           │
                    ┌──────────────────────┼──────────────────────┐
                    ▼                      ▼                      ▼
            ┌─────────────┐       ┌─────────────┐       ┌─────────────┐
            │Inventory Svc│       │Payment Svc  │       │Notification │
            │   (gRPC)    │       │   (gRPC)    │       │   (NATS)    │
            └─────────────┘       └─────────────┘       └─────────────┘
```

## Services

### 1. User Service

Handles user authentication, registration, and profile management.

**Responsibilities:**

- User registration and login
- JWT token generation and validation
- Profile CRUD operations
- Password hashing (bcrypt)
- Email verification

**Endpoints:**

```
POST   /api/v1/users/register
POST   /api/v1/users/login
GET    /api/v1/users/profile
PUT    /api/v1/users/profile
POST   /api/v1/users/verify-email
POST   /api/v1/users/refresh-token
```

**Database Schema:**

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    phone VARCHAR(20),
    email_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

### 2. Order Service

Manages order lifecycle from creation to fulfillment.

**Responsibilities:**

- Order creation and management
- Order status tracking
- Order history
- Event publishing for downstream services

**Events Published:**

- `order.created`
- `order.paid`
- `order.shipped`
- `order.delivered`
- `order.cancelled`

### 3. Payment Service

Handles payment processing with multiple providers.

**Responsibilities:**

- Payment processing
- Refund handling
- Payment method storage (PCI compliant)
- Webhook handling

**Supported Providers:**

- Stripe
- PayPal
- Square

### 4. Inventory Service

Manages product inventory and stock levels.

**Responsibilities:**

- Stock level tracking
- Reservation system
- Low stock alerts
- Inventory reconciliation

### 5. Notification Service

Sends notifications across multiple channels.

**Responsibilities:**

- Email notifications (SendGrid/AWS SES)
- SMS notifications (Twilio)
- Push notifications (Firebase)
- In-app notifications

## Getting Started

### Prerequisites

- Go 1.21 or later
- Docker 24.0+ and Docker Compose
- Kubernetes 1.28+ (optional, for k8s deployment)
- kubectl and Helm 3
- make

### Quick Start with Docker Compose

```bash
# Clone the repository
git clone https://github.com/example/microservices-platform
cd microservices-platform

# Start all services
make docker-up

# Or manually
docker-compose up -d

# Verify services are running
docker-compose ps

# View logs
docker-compose logs -f [service-name]
```

### Accessing Services

| Service | URL | Credentials |
|---------|-----|-------------|
| API Gateway | <http://localhost:8080> | - |
| User Service | <http://localhost:8081> | - |
| Order Service | <http://localhost:8082> | - |
| Payment Service | <http://localhost:8083> | - |
| Inventory Service | <http://localhost:8084> | - |
| PostgreSQL | localhost:5432 | postgres/postgres |
| MongoDB | localhost:27017 | - |
| Redis | localhost:6379 | - |
| Kafka UI | <http://localhost:8090> | - |
| Grafana | <http://localhost:3000> | admin/admin |
| Jaeger | <http://localhost:16686> | - |

### Running Individual Services

```bash
# User Service
cd services/user-service
go mod download
go run cmd/main.go

# Order Service
cd services/order-service
go mod download
go run cmd/main.go
```

## Development Guide

### Project Structure

```
microservices-platform/
├── services/
│   ├── user-service/
│   │   ├── cmd/
│   │   │   └── main.go
│   │   ├── internal/
│   │   │   ├── config/
│   │   │   ├── domain/
│   │   │   ├── repository/
│   │   │   ├── service/
│   │   │   ├── handler/
│   │   │   └── middleware/
│   │   ├── api/
│   │   │   └── proto/
│   │   ├── migrations/
│   │   ├── tests/
│   │   ├── Dockerfile
│   │   ├── go.mod
│   │   └── go.sum
│   ├── order-service/
│   ├── payment-service/
│   ├── inventory-service/
│   ├── notification-service/
│   ├── shipping-service/
│   └── analytics-service/
├── api-gateway/
│   ├── kong/
│   └── envoy/
├── infra/
│   ├── docker/
│   ├── k8s/
│   ├── terraform/
│   └── helm/
├── shared/
│   ├── pkg/
│   │   ├── logger/
│   │   ├── tracer/
│   │   ├── metrics/
│   │   ├── errors/
│   │   ├── middleware/
│   │   └── utils/
│   └── proto/
├── deployments/
│   ├── docker-compose.yml
│   ├── docker-compose.prod.yml
│   └── Makefile
├── docs/
│   ├── architecture/
│   ├── api/
│   └── runbooks/
└── scripts/
    ├── setup.sh
    ├── test.sh
    └── deploy.sh
```

### Coding Standards

1. **Package Structure**: Follow Clean Architecture / Domain-Driven Design
2. **Error Handling**: Use custom error types with error codes
3. **Logging**: Structured JSON logging with correlation IDs
4. **Testing**: Unit tests (>80% coverage), integration tests, e2e tests
5. **Documentation**: Go doc comments, OpenAPI specs
6. **Code Quality**: golangci-lint, go fmt, go vet

### Creating a New Service

```bash
# Use the service generator
make new-service SERVICE_NAME=review-service

# Or manually
cd services
cp -r user-service review-service
cd review-service
# Update go.mod, Dockerfile, configs
```

## Deployment

### Docker Compose Deployment

```bash
# Development
docker-compose -f docker-compose.yml up -d

# Production
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d

# Scale specific services
docker-compose up -d --scale order-service=3
```

### Kubernetes Deployment

```bash
# Create namespace
kubectl create namespace microservices

# Deploy using Helm
helm install microservices ./helm/microservices \
  --namespace microservices \
  --values ./helm/values.yaml

# Or apply manifests directly
kubectl apply -f infra/k8s/namespace.yaml
kubectl apply -f infra/k8s/configmaps/
kubectl apply -f infra/k8s/secrets/
kubectl apply -f infra/k8s/services/
kubectl apply -f infra/k8s/deployments/

# Verify deployment
kubectl get pods -n microservices
kubectl get svc -n microservices
```

### Istio Service Mesh Setup

```bash
# Install Istio
istioctl install --set profile=default -y

# Enable sidecar injection
kubectl label namespace microservices istio-injection=enabled

# Apply Istio configurations
kubectl apply -f infra/istio/gateway.yaml
kubectl apply -f infra/istio/virtualservices.yaml
kubectl apply -f infra/istio/destinationrules.yaml
kubectl apply -f infra/istio/policies.yaml

# View mesh dashboard
istioctl dashboard kiali
```

### CI/CD Pipeline

```yaml
# .github/workflows/deploy.yml
name: Deploy Microservices

on:
  push:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - run: make test

  build:
    needs: test
    runs-on: ubuntu-latest
    strategy:
      matrix:
        service: [user, order, payment, inventory, notification]
    steps:
      - uses: actions/checkout@v3
      - uses: docker/login-action@v2
        with:
          username: ${{ secrets.DOCKER_USERNAME }}
          password: ${{ secrets.DOCKER_PASSWORD }}
      - run: |
          docker build -t ${{ secrets.DOCKER_USERNAME }}/${{ matrix.service }}-service:${{ github.sha }} \
            -f services/${{ matrix.service }}-service/Dockerfile .
          docker push ${{ secrets.DOCKER_USERNAME }}/${{ matrix.service }}-service:${{ github.sha }}

  deploy:
    needs: build
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Deploy to Kubernetes
        run: |
          kubectl set image deployment/user-service \
            user-service=${{ secrets.DOCKER_USERNAME }}/user-service:${{ github.sha }}
```

## Observability

### Distributed Tracing with Jaeger

All services are instrumented with OpenTelemetry for distributed tracing.

```go
// Example trace span
tracer := otel.Tracer("user-service")
ctx, span := tracer.Start(ctx, "CreateUser")
defer span.End()

span.SetAttributes(
    attribute.String("user.email", email),
    attribute.Int64("user.id", userID),
)
```

Access Jaeger UI: <http://localhost:16686>

### Metrics with Prometheus

Key metrics exposed:

- HTTP request duration/latency
- Request rate and error rate
- Database connection pool stats
- Cache hit/miss ratios
- Message queue lag

```promql
# Request rate by service
rate(http_requests_total[5m])

# 95th percentile latency
histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))

# Error rate
rate(http_requests_total{status=~"5.."}[5m])
```

### Logging with Loki

Structured JSON logging with correlation IDs:

```json
{
  "level": "info",
  "timestamp": "2024-01-15T10:30:00Z",
  "service": "user-service",
  "trace_id": "abc123def456",
  "span_id": "xyz789",
  "message": "User created successfully",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "duration_ms": 45
}
```

### Dashboards

Pre-configured Grafana dashboards:

- Service Overview
- Request Metrics
- Database Performance
- Cache Performance
- Infrastructure Metrics

## Security

### Authentication & Authorization

- JWT tokens with RS256 signing
- Refresh token rotation
- OAuth 2.0 / OpenID Connect support
- RBAC with Casbin

### mTLS

Istio automatically enables mTLS between services:

```yaml
apiVersion: security.istio.io/v1beta1
kind: PeerAuthentication
metadata:
  name: default
  namespace: microservices
spec:
  mtls:
    mode: STRICT
```

### Secrets Management

Using HashiCorp Vault or Kubernetes Secrets:

```bash
# Store secret in Vault
vault kv put secret/user-service \
  db_password=secure_password \
  jwt_secret=super_secret_key

# Or create K8s secret
kubectl create secret generic user-service-secrets \
  --from-literal=db_password=secure_password
```

### Security Scanning

```bash
# Container vulnerability scan
trivy image user-service:latest

# Dependency vulnerability scan
govulncheck ./...

# SAST scanning
semgrep --config=auto .
```

## Performance

### Load Testing

Using k6 for load testing:

```javascript
// tests/load/order_creation.js
import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '2m', target: 100 },
    { duration: '5m', target: 500 },
    { duration: '2m', target: 1000 },
    { duration: '2m', target: 0 },
  ],
};

export default function () {
  const payload = JSON.stringify({
    items: [{ product_id: '123', quantity: 2 }],
    shipping_address: { ... }
  });

  const res = http.post('http://localhost:8080/api/v1/orders', payload, {
    headers: { 'Content-Type': 'application/json' },
  });

  check(res, {
    'status is 201': (r) => r.status === 201,
    'response time < 500ms': (r) => r.timings.duration < 500,
  });

  sleep(1);
}
```

Run load test:

```bash
k6 run tests/load/order_creation.js
```

### Benchmark Results

| Metric | Target | Actual |
|--------|--------|--------|
| P50 Latency | < 50ms | 35ms |
| P95 Latency | < 200ms | 150ms |
| P99 Latency | < 500ms | 380ms |
| Throughput | 10,000 RPS | 12,000 RPS |
| Error Rate | < 0.1% | 0.05% |
| Availability | 99.9% | 99.95% |

### Optimization Techniques

1. **Caching Strategy**
   - Redis for hot data
   - Application-level caching
   - CDN for static assets

2. **Database Optimization**
   - Connection pooling
   - Read replicas
   - Query optimization
   - Proper indexing

3. **Async Processing**
   - Event-driven architecture
   - Background job processing
   - Message queues for decoupling

4. **Resource Management**
   - Horizontal Pod Autoscaling
   - Vertical Pod Autoscaling
   - Cluster Autoscaling

## Troubleshooting

### Common Issues

#### Service Discovery Failures

```bash
# Check Consul health
kubectl exec -it consul-0 -- consul members

# Check DNS resolution
kubectl run -it --rm debug --image=busybox:1.28 --restart=Never -- nslookup user-service
```

#### Database Connection Issues

```bash
# Check PostgreSQL logs
kubectl logs -l app=postgres --tail=100

# Test connection
psql -h localhost -U postgres -d users -c "SELECT 1"
```

#### Message Queue Lag

```bash
# Check Kafka consumer lag
kubectl exec -it kafka-0 -- kafka-consumer-groups.sh \
  --bootstrap-server localhost:9092 \
  --describe \
  --group order-consumers
```

#### Memory Leaks

```bash
# Get heap profile
curl http://user-service:8080/debug/pprof/heap > heap.prof
go tool pprof heap.prof

# View goroutines
curl http://user-service:8080/debug/pprof/goroutine?debug=1
```

### Debugging Commands

```bash
# Port forward to service
kubectl port-forward svc/user-service 8081:80

# Execute into container
kubectl exec -it deployment/user-service -- /bin/sh

# Check resource usage
kubectl top pods -n microservices

# View events
kubectl get events -n microservices --sort-by=.lastTimestamp
```

### Support

- Documentation: <https://docs.example.com/microservices>
- Issues: <https://github.com/example/microservices-platform/issues>
- Slack: #microservices-support

## License

MIT License - See LICENSE file for details

## Contributing

Please read CONTRIBUTING.md for details on our code of conduct and the process for submitting pull requests.

---

**Maintained by**: Platform Engineering Team
**Last Updated**: 2024-01-15
**Version**: 2.1.0
