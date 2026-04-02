# 部署指南

> **版本**: v1.0  
> **更新日期**: 2026-04-02

---

## 🚀 快速部署

### Docker Compose (推荐)

```bash
# 1. 配置环境变量
cp .env.example .env
# 编辑 .env 文件

# 2. 启动所有服务
docker-compose up -d

# 3. 检查状态
docker-compose ps
```

---

## 📋 环境变量

| 变量 | 说明 | 必需 |
|------|------|------|
| `DB_HOST` | 数据库主机 | ✅ |
| `DB_USER` | 数据库用户 | ✅ |
| `DB_PASSWORD` | 数据库密码 | ✅ |
| `JWT_SECRET` | JWT 签名密钥 | ✅ |

---

## ☸️ Kubernetes

```bash
kubectl apply -f k8s/
```

---

*最后更新: 2026-04-02*
