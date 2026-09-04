# TS-017: Grafana Dashboard Design - Visualization Best Practices

> **维度**: Technology Stack
> **级别**: S (36 KB)
> **标签**: #grafana #dashboard #visualization #observability #monitoring
> **权威来源**:
>
> - [Grafana Documentation](https://grafana.com/docs/) - Grafana Labs
> - [Dashboard Best Practices](https://grafana.com/docs/grafana/latest/dashboards/build-dashboards/best-practices/) - Grafana Docs
> - [Grafana Academy](https://grafana.com/academy/) - Grafana Labs

> **Go 版本**: 1.27+
---

## 1. Grafana Architecture

### 1.1 System Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Grafana System Architecture                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Grafana Frontend                                    │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │  React/TypeScript Application                                    │  │  │
│  │  │                                                                  │  │  │
│  │  │  Components:                                                     │  │  │
│  │  │  • Panel Renderer (Graph, Table, Stat, Gauge, Heatmap, etc.)    │  │  │
│  │  │  • Dashboard Grid (react-grid-layout)                           │  │  │
│  │  │  • Query Editor (per data source)                               │  │  │
│  │  │  • Variable Selector                                            │  │  │
│  │  │  • Alert Rule Editor                                            │  │  │
│  │  │                                                                  │  │  │
│  │  │  State Management: Redux + Redux Toolkit                        │  │  │
│  │  │                                                                  │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                 │                                            │
│                                 │ HTTP/WebSocket                             │
│                                 ▼                                            │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Grafana Backend (Go)                                │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │  API Layer (HTTP Server)                                        │  │  │
│  │  │                                                                  │  │  │
│  │  │  Routes:                                                        │  │  │
│  │  │  • /api/dashboards (CRUD dashboards)                            │  │  │
│  │  │  • /api/datasources (Data source management)                    │  │  │
│  │  │  • /api/annotations (Event annotations)                         │  │  │
│  │  │  • /api/alert-rules (Alerting)                                  │  │  │
│  │  │  • /api/live (WebSocket streaming)                              │  │  │
│  │  │                                                                  │  │  │
│  │  │  Middleware:                                                    │  │  │
│  │  │  • Authentication (JWT, OAuth, SAML, LDAP)                      │  │  │
│  │  │  • Authorization (RBAC)                                         │  │  │
│  │  │  • Rate Limiting                                                │  │  │
│  │  │                                                                  │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                  │                                     │  │
│  │                                  ▼                                     │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │  Data Source Proxy / Plugin System                               │  │  │
│  │  │                                                                  │  │  │
│  │  │  Plugin Architecture:                                            │  │  │
│  │  │  • Core data sources (Prometheus, InfluxDB, Elasticsearch, etc.)│  │  │
│  │  │  • External plugins (loaded dynamically)                        │  │  │
│  │  │  • Backend plugins (Go SDK)                                     │  │  │
│  │  │                                                                  │  │  │
│  │  │  Query Processing:                                              │  │  │
│  │  │  1. Frontend sends query request                                │  │  │
│  │  │  2. Backend routes to appropriate data source                   │  │  │
│  │  │  3. Data source plugin translates to native query               │  │  │
│  │  │  4. Execute against external system                             │  │  │
│  │  │  5. Transform response to Grafana data frames                   │  │  │
│  │  │  6. Return to frontend                                          │  │  │
│  │  │                                                                  │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                  │                                     │  │
│  │                                  ▼                                     │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │  Storage Layer                                                   │  │  │
│  │  │                                                                  │  │  │
│  │  │  SQLite (default) / MySQL / PostgreSQL:                         │  │  │
│  │  │  • Dashboard definitions (JSON)                                 │  │  │
│  │  │  • User/Organization data                                       │  │  │
│  │  │  • Data source configurations                                   │  │  │
│  │  │  • Alert rules and notifications                                │  │  │
│  │  │                                                                  │  │  │
│  │  │  Optional External Storage:                                     │  │  │
│  │  │  • S3/GCS for image storage                                     │  │  │
│  │  │  • Redis for session/cache                                      │  │  │
│  │  │                                                                  │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                 │                                            │
│                                 │ Query                                      │
│                                 ▼                                            │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Data Sources                                        │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  │  │
│  │  │ Prometheus  │  │ InfluxDB    │  │ Elasticsearch│  │ Loki        │  │  │
│  │  │ (Metrics)   │  │ (Time Series)│  │ (Logs)      │  │ (Logs)      │  │  │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘  │  │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  │  │
│  │  │ MySQL       │  │ PostgreSQL  │  │ CloudWatch  │  │ Azure Monitor│  │  │
│  │  │ (SQL)       │  │ (SQL)       │  │ (AWS)       │  │ (Azure)     │  │  │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘  │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Dashboard Design Patterns

### 2.1 The RED Method Dashboard

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    RED Method Dashboard Layout                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │  Service: API Gateway    Time Range: Last 1 Hour    Refresh: 10s      │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│  ┌───────────────────────────────┐  ┌───────────────────────────────┐      │
│  │  REQUEST RATE (R)             │  │  ERROR RATE (E)               │      │
│  │  ┌─────────────────────────┐  │  │  ┌─────────────────────────┐  │      │
│  │  │                         │  │  │  │                         │  │      │
│  │  │      ╱╲      ╱╲        │  │  │  │     ╱                   │  │      │
│  │  │     ╱  ╲    ╱  ╲╱╲     │  │  │  │    ╱                    │  │      │
│  │  │    ╱    ╲  ╱      ╲    │  │  │  │   ╱                     │  │      │
│  │  │   ╱      ╲╱            │  │  │  │  ╱                      │  │      │
│  │  │  ╱                      │  │  │  │ ╱                       │  │      │
│  │  └─────────────────────────┘  │  │  └─────────────────────────┘  │      │
│  │  Current: 1.2k rps            │  │  Current: 0.5%              │      │
│  │  Target: < 2k rps             │  │  Target: < 1%               │      │
│  └───────────────────────────────┘  └───────────────────────────────┘      │
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │  DURATION (D) - P99 Latency                                           │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │                                                                 │  │  │
│  │  │     ╱╲      ╱╲                                                  │  │  │
│  │  │    ╱  ╲    ╱  ╲╱╲     ─── Threshold: 500ms                     │  │  │
│  │  │   ╱    ╲  ╱      ╲                                            │  │  │
│  │  │  ╱      ╲╱                                                    │  │  │
│  │  │                                                                 │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │  Current P99: 320ms | P95: 180ms | P50: 50ms                          │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│  ┌───────────────────────────────┐  ┌───────────────────────────────┐      │
│  │  Top 5 Slowest Endpoints      │  │  Error Rate by Status Code    │      │
│  │  ┌─────────────────────────┐  │  │  ┌─────────────────────────┐  │      │
│  │  │ /api/v1/orders   450ms  │  │  │  │ 500: 0.3%               │  │      │
│  │  │ /api/v1/users    320ms  │  │  │  │ 502: 0.1%               │  │      │
│  │  │ /api/v1/products 180ms  │  │  │  │ 503: 0.05%              │  │      │
│  │  │ /health          15ms   │  │  │  │ 504: 0.05%              │  │      │
│  │  │ /metrics         5ms    │  │  │  │                         │  │      │
│  │  └─────────────────────────┘  │  │  └─────────────────────────┘  │      │
│  └───────────────────────────────┘  └───────────────────────────────┘      │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 2.2 The USE Method Dashboard

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    USE Method Dashboard Layout                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │  Resource: Database Servers    Time Range: Last 6 Hours               │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │  UTILIZATION                                                          │  │
│  │  ┌──────────────────┐ ┌──────────────────┐ ┌──────────────────┐      │  │
│  │  │ CPU Usage        │ │ Memory Usage     │ │ Disk Usage       │      │  │
│  │  │ ┌──────────────┐ │ │ ┌──────────────┐ │ │ ┌──────────────┐ │      │  │
│  │  │ │████████░░░░░│ │ │ │█████████░░░░░│ │ │ │██████░░░░░░░░│ │      │  │
│  │  │ │   72%      │ │ │ │   85%      │ │ │ │   45%      │ │      │  │
│  │  │ └──────────────┘ │ │ └──────────────┘ │ │ └──────────────┘ │      │  │
│  │  │ Alert: > 80%     │ │ Alert: > 90%     │ │ Alert: > 85%     │      │  │
│  │  └──────────────────┘ └──────────────────┘ └──────────────────┘      │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │  SATURATION                                                           │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │  Load Average vs CPU Count                                      │  │  │
│  │  │                                                                 │  │  │
│  │  │  Load  ████████████████████                                     │  │  │
│  │  │  CPUs  ████████                                                 │  │  │
│  │  │                                                                 │  │  │
│  │  │  Ratio: 3.2 (load > 2x CPUs = saturated)                       │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │  Disk I/O Saturation (% of time with I/O in progress)           │  │  │
│  │  │                                                                 │  │  │
│  │  │  ████████░░░░░░░░░░░░  35%                                    │  │  │
│  │  │                                                                 │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │  ERRORS                                                               │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │  Disk I/O Errors                                                │  │  │
│  │  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │  │  │
│  │  │  │ db-01: 0     │  │ db-02: 2     │  │ db-03: 0     │          │  │  │
│  │  │  └──────────────┘  └──────────────┘  └──────────────┘          │  │  │
│  │  │                                                                 │  │  │
│  │  │  Network Errors                                                 │  │  │
│  │  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │  │  │
│  │  │  │ db-01: 0     │  │ db-02: 0     │  │ db-03: 1     │          │  │  │
│  │  │  └──────────────┘  └──────────────┘  └──────────────┘          │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 3. Dashboard as Code (JSON Model)

```json
{
  "dashboard": {
    "title": "API Gateway Overview",
    "tags": ["api", "gateway", "production"],
    "timezone": "browser",
    "schemaVersion": 36,
    "refresh": "10s",
    "time": {
      "from": "now-1h",
      "to": "now"
    },
    "templating": {
      "list": [
        {
          "name": "service",
          "type": "query",
          "query": "label_values(http_requests_total, service)",
          "multi": true,
          "includeAll": true
        },
        {
          "name": "env",
          "type": "custom",
          "query": "production,staging,development",
          "current": {"text": "production", "value": "production"}
        }
      ]
    },
    "panels": [
      {
        "id": 1,
        "title": "Request Rate",
        "type": "timeseries",
        "targets": [
          {
            "expr": "sum(rate(http_requests_total{service=~\"$service\",env=\"$env\"}[5m]))",
            "legendFormat": "{{service}}"
          }
        ],
        "fieldConfig": {
          "defaults": {
            "unit": "reqps",
            "min": 0,
            "custom": {
              "drawStyle": "line",
              "lineWidth": 2,
              "fillOpacity": 10
            }
          }
        }
      },
      {
        "id": 2,
        "title": "Error Rate",
        "type": "stat",
        "targets": [
          {
            "expr": "sum(rate(http_requests_total{service=~\"$service\",env=\"$env\",status=~\"5..\"}[5m])) / sum(rate(http_requests_total{service=~\"$service\",env=\"$env\"}[5m]))",
            "legendFormat": "Error %"
          }
        ],
        "fieldConfig": {
          "defaults": {
            "unit": "percentunit",
            "thresholds": {
              "steps": [
                {"color": "green", "value": 0},
                {"color": "yellow", "value": 0.01},
                {"color": "red", "value": 0.05}
              ]
            }
          }
        }
      }
    ]
  }
}
```

---

## 4. Configuration Best Practices

```ini
# grafana.ini
[server]
protocol = http
http_port = 3000
domain = grafana.example.com
root_url = https://grafana.example.com

[database]
type = postgres
host = localhost:5432
name = grafana
user = grafana
password = $__env{GF_DB_PASSWORD}

[security]
admin_user = admin
admin_password = $__env{GF_SECURITY_ADMIN_PASSWORD}
secret_key = $__env{GF_SECURITY_SECRET_KEY}

[auth]
disable_login_form = false
disable_signout_menu = false

[auth.generic_oauth]
enabled = true
name = SSO
allow_sign_up = true
client_id = $__env{GF_AUTH_GENERIC_OAUTH_CLIENT_ID}
client_secret = $__env{GF_AUTH_GENERIC_OAUTH_CLIENT_SECRET
scopes = openid profile email
token_url = https://auth.example.com/oauth/token
api_url = https://auth.example.com/oauth/userinfo

[analytics]
reporting_enabled = false
check_for_updates = true

[unified_alerting]
enabled = true
execute_alerts = true
```

---

## 5. Visual Representations

### Panel Type Selection Guide

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Grafana Panel Type Selection Guide                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Time Series Data (Metrics)                                                │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │  • Time Series (new) - Modern, preferred for metrics                  │  │
│  │  • Graph (legacy) - Old panel, migrate to Time Series                │  │
│  │  • Heatmap - Show distribution over time                             │  │
│  │  • Histogram - Show value distribution                               │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│  Single Values (Gauges)                                                    │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │  • Stat - Big number with sparkline (CPU, memory usage)              │  │
│  │  • Gauge - Circular gauge with thresholds                            │  │
│  │  • Bar Gauge - Horizontal/vertical bar with limits                   │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│  Geospatial Data                                                           │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │  • Geomap - Map visualization for location data                      │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│  Tabular Data                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │  • Table - Sortable, filterable table                                │  │
│  │  • Logs - Formatted log viewing with highlighting                    │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│  Special Purpose                                                           │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │  • Node Graph - Network topology visualization                       │  │
│  │  • Canvas - Custom visualizations                                    │  │
│  │  • Alert List - Show active alerts                                   │  │
│  │  • Annotations List - Show events                                    │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 6. References

1. **Grafana Documentation** (2024). grafana.com/docs
2. **Grafana Academy** (2024). grafana.com/academy
3. **Wilkerson, M.** (2021). Grafana 8 Fundamentals. Packt Publishing.

---

*Document Version: 1.0 | Last Updated: 2024*
