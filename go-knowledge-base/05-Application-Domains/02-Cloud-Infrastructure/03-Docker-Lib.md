# Docker SDK

> **分类**: 成熟应用领域

---

## 安装

```go
import "github.com/docker/docker/client"
```

---

## 连接

```go
cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
if err != nil {
    log.Fatal(err)
}
```

---

## 容器操作

### 列出容器

```go
containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
if err != nil {
    log.Fatal(err)
}

for _, container := range containers {
    fmt.Println(container.ID, container.Image)
}
```

### 创建并启动

```go
resp, err := cli.ContainerCreate(ctx, &container.Config{
    Image: "alpine",
    Cmd:   []string{"echo", "hello world"},
}, nil, nil, nil, "")

if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
    log.Fatal(err)
}
```

---

## 镜像操作

```go
// 拉取镜像
reader, err := cli.ImagePull(ctx, "alpine:latest", types.ImagePullOptions{}),
io.Copy(os.Stdout, reader)
```
