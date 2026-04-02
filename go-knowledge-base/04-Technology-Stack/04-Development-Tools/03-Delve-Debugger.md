# Delve 调试器

> **分类**: 开源技术堆栈

---

## 安装

```bash
go install github.com/go-delve/delve/cmd/dlv@latest
```

---

## 基本命令

```bash
# 调试当前目录程序
dlv debug

# 调试指定包
dlv debug github.com/user/project

# 附加到运行进程
dlv attach <pid>

# 调试测试
dlv test
```

---

## 常用命令

```
(dlv) break main.main        # 设置断点
(dlv) continue               # 继续执行
(dlv) next                   # 单步跳过
(dlv) step                   # 单步进入
(dlv) stepout                # 跳出函数
(dlv) print variable         # 打印变量
(dlv) locals                 # 显示局部变量
(dlv) args                   # 显示参数
(dlv) stack                  # 显示调用栈
(dlv) goroutines             # 显示 goroutines
(dlv) quit                   # 退出
```

---

## VS Code 集成

### launch.json

```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${fileDirname}",
            "env": {},
            "args": []
        }
    ]
}
```

---

## 远程调试

```bash
# 服务器
dlv debug --headless --listen=:2345 --api-version=2

# 本地连接
dlv connect server:2345
```
