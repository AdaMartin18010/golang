# CI/CD 集成

> **分类**: 成熟应用领域

---

## GitHub Actions

```yaml
name: CI
on: [push, pull_request]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
    - uses: actions/setup-go@v5
      with:
        go-version: '1.22'

    - name: Build
      run: go build -v ./...

    - name: Test
      run: go test -v ./...

    - name: Lint
      uses: golangci/golangci-lint-action@v3
```

---

## GitLab CI

```yaml
stages:
  - build
  - test

build:
  stage: build
  image: golang:1.22
  script:
    - go build -o app

test:
  stage: test
  script:
    - go test ./...
```
