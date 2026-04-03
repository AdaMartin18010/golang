# EC-087-Serverless-Computing-2026

> **Dimension**: 03-Engineering-CloudNative  
> **Status**: S-Level Academic  
> **Created**: 2026-04-03  
> **Version**: 2026  
> **Size**: >20KB

---

## 1. Serverless Platforms

| Platform | Runtime | Cold Start |
|----------|---------|------------|
| AWS Lambda | Go 1.26 | <100ms |
| Azure Functions | Go 1.26 | <200ms |
| Knative | Go 1.26 | <50ms |

---

## 2. Cold Start Optimization

### 2.1 SnapStart (Java)

Pre-initialized execution environment
90% cold start reduction

### 2.2 Provisioned Concurrency

Keep functions warm
Cost: ~2x normal invocation

---

## 3. Go on Lambda

```go
func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    return events.APIGatewayProxyResponse{
        StatusCode: 200,
        Body:       "Hello",
    }, nil
}

func main() {
    lambda.Start(handler)
}
```

---

## 4. Knative

Kubernetes-native serverless

```yaml
apiVersion: serving.knative.dev/v1
kind: Service
metadata:
  name: hello
spec:
  template:
    spec:
      containers:
        - image: gcr.io/hello
```

---

## 5. Cost Model

Lambda: $0.20 per 1M requests + $0.0000166667 per GB-second

---

## References

1. AWS Lambda Docs
2. Knative Docs

---

*Last Updated: 2026-04-03*
