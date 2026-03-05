# 快速启动指南

> 5分钟内启动并运行基础环境

---

## 前置要求

- Docker & Docker Compose
- Go 1.21+
- Git

## 一键启动

```bash
# 1. 克隆本仓库
git clone <repo-url>
cd ai-native-observability-analysis

# 2. 启动基础设施
docker-compose up -d

# 3. 验证安装
docker-compose ps
```

预期输出：

```
        Name                      Command               State            Ports
-----------------------------------------------------------------------------------------
observability-analyzer-neo4j   /startup/docker-entrypoint.sh    Up      7474/tcp, 7687/tcp
observability-analyzer-weaviate  /bin/weaviate --host 0.0.0.0   Up      8080/tcp
```

## 第一步：分析示例项目

```bash
# 下载分析目标（以OTel Collector为例）
mkdir -p targets
cd targets
git clone --depth 1 https://github.com/open-telemetry/opentelemetry-collector.git
cd ..

# 运行分析器
go run ./cmd/analyze \
  --project=./targets/opentelemetry-collector \
  --output=./data/otelcol

# 验证输出
ls -la ./data/otelcol/
# 应包含：components.json, interfaces.json, configs.json, metrics.json
```

## 第二步：导入知识图谱

```bash
# 导入到Neo4j
go run ./cmd/import \
  --input=./data/otelcol \
  --neo4j-uri=bolt://localhost:7687 \
  --neo4j-user=neo4j \
  --neo4j-password=password

# 验证导入
open http://localhost:7474
# 用户名: neo4j
# 密码: password
# 运行查询: MATCH (n) RETURN count(n)
```

## 第三步：启动MCP服务

```bash
# 启动MCP服务器
go run ./cmd/mcp-server \
  --neo4j-uri=bolt://localhost:7687 \
  --port=8080

# 服务将在 http://localhost:8080/mcp 可用
```

## 第四步：测试查询

### 使用curl

```bash
# 查询项目信息
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "tool": "query_project",
    "args": {
      "name": "OpenTelemetry Collector"
    }
  }'

# 搜索组件
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "tool": "search_components",
    "args": {
      "type": "processor",
      "capability": "batch"
    }
  }'
```

### 使用Python客户端

```bash
# 安装客户端
pip install mcp-client
```

```python
# test_client.py
import mcp

client = mcp.Client("http://localhost:8080/mcp")

# 查询项目
result = client.call("query_project", {
    "name": "OpenTelemetry Collector"
})
print(result)

# 获取组件详情
result = client.call("get_component", {
    "id": "otelcol-processor-batch"
})
print(result)

# 解释指标
result = client.call("explain_metric", {
    "metric": "processor_batch_batch_send_size"
})
print(result)
```

```bash
python test_client.py
```

## 第五步：使用自然语言查询

```bash
# 启动NL查询服务（需要OpenAI API Key）
export OPENAI_API_KEY=your-key

go run ./cmd/nl-query \
  --neo4j-uri=bolt://localhost:7687 \
  --port=8081

# 测试查询
curl -X POST http://localhost:8081/query \
  -H "Content-Type: application/json" \
  -d '{"query": "如何配置batch处理器来提高吞吐量？"}'
```

## 常见问题

### Neo4j连接失败

```bash
# 检查容器状态
docker-compose logs neo4j

# 重置Neo4j
docker-compose down -v
docker-compose up -d neo4j
```

### 分析器找不到Go文件

```bash
# 确保目标目录正确
ls targets/opentelemetry-collector/processor/

# 检查Go版本
go version
```

### MCP服务无响应

```bash
# 检查端口占用
lsof -i :8080

# 查看日志
go run ./cmd/mcp-server 2>&1 | tee mcp.log
```

## 下一步

- 📖 阅读 [OTel Collector分析](./01-OTEL-COLLECTOR-ANALYSIS.md) 了解深度分析方法论
- 🔧 查看 [实现路线图](./03-IMPLEMENTATION-ROADMAP.md) 了解完整开发计划
- 🏗️ 尝试分析你自己的Go项目

## 获取帮助

- 查看完整文档: [README.md](./README.md)
- 提交Issue: [GitHub Issues]
- 加入讨论: [Discord/Slack]

---

*遇到问题？请查看 [故障排除指南](./TROUBLESHOOTING.md)*
