# 构建自动化

> **分类**: 成熟应用领域

---

## goreleaser

```yaml
# .goreleaser.yaml
builds:
  - binary: myapp
    goos:
      - linux
      - darwin
      - windows
    goarch:
      - amd64
      - arm64

dockers:
  - image_templates:
      - "myapp:{{ .Tag }}"
```

---

## 使用

```bash
# 发布
goreleaser release

# 快照构建
goreleaser release --snapshot --clean
```

---

## 交叉编译

```bash
GOOS=linux GOARCH=amd64 go build -o bin/linux/myapp
GOOS=darwin GOARCH=arm64 go build -o bin/macos/myapp
GOOS=windows GOARCH=amd64 go build -o bin/windows/myapp.exe
```
