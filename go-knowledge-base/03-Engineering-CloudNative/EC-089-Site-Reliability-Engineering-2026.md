# EC-089-Site-Reliability-Engineering-2026

> **Dimension**: 03-Engineering-CloudNative  
> **Status**: S-Level Academic  
> **Created**: 2026-04-03  
> **Version**: SRE 2026  
> **Size**: >20KB

---

## 1. SLI/SLO/SLA

| Term | Definition | Example |
|------|------------|---------|
| SLI | Indicator | Latency, error rate |
| SLO | Objective | 99.9% availability |
| SLA | Agreement | Refund if SLO missed |

---

## 2. Error Budget

```
Error Budget = (1 - SLO) * time window

Example:
SLO = 99.9%
Window = 30 days
Error Budget = 0.1% * 30 days = 43 minutes
```

---

## 3. Four Golden Signals

1. Latency
2. Traffic
3. Errors
4. Saturation

---

## 4. Chaos Engineering

Principles:
1. Define steady state
2. Hypothesize
3. Inject failure
4. Measure

---

## 5. Tools

- Prometheus
- Grafana
- Jaeger
- Chaos Mesh

---

## References

1. Google SRE Book
2. Chaos Engineering Book

---

*Last Updated: 2026-04-03*
