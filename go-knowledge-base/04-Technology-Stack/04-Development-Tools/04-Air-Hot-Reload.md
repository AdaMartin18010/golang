# Air 热重载

> **分类**: 开源技术堆栈

---

## 安装

```bash
go install github.com/cosmtrek/air@latest
```

---

## 配置

### .air.toml

```toml
root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
  bin = "./tmp/main"
  cmd = "go build -o ./tmp/main ."
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor"]
  include_ext = ["go", "tpl", "tmpl", "html"]
  exclude_regex = ["_test.go"]

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[misc]
  clean_on_exit = false
```

---

## 使用

```bash
# 初始化配置
air init

# 运行
air

# 指定配置
air -c .air.toml
```

---

## Docker 使用

```dockerfile
FROM cosmtrek/air:latest

WORKDIR /app
COPY . .

CMD ["air", "-c", ".air.toml"]
```
