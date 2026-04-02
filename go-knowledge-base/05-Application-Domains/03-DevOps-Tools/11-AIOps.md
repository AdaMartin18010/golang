# AIOps 基础

> **分类**: 成熟应用领域
> **标签**: #aiops #mlops #observability

---

## 异常检测

### 基于统计的检测

```go
type AnomalyDetector struct {
    windowSize int
    threshold  float64
    history    []float64
    mu         sync.RWMutex
}

func (d *AnomalyDetector) Update(value float64) bool {
    d.mu.Lock()
    defer d.mu.Unlock()

    d.history = append(d.history, value)
    if len(d.history) > d.windowSize {
        d.history = d.history[1:]
    }

    if len(d.history) < d.windowSize {
        return false
    }

    mean, std := d.calculateStats()
    zScore := math.Abs(value-mean) / std

    return zScore > d.threshold
}

func (d *AnomalyDetector) calculateStats() (mean, std float64) {
    sum := 0.0
    for _, v := range d.history {
        sum += v
    }
    mean = sum / float64(len(d.history))

    variance := 0.0
    for _, v := range d.history {
        variance += math.Pow(v-mean, 2)
    }
    std = math.Sqrt(variance / float64(len(d.history)))

    return
}
```

---

## 日志聚类

```go
// 日志模板提取
func ExtractTemplate(log string) string {
    // 移除变量（数字、IP、UUID等）
    patterns := []string{
        `\d+`,                    // 数字
        `\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`,  // IP
        `[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`,  // UUID
    }

    result := log
    for _, pattern := range patterns {
        re := regexp.MustCompile(pattern)
        result = re.ReplaceAllString(result, "*")
    }

    return result
}
```

---

## 预测性扩容

```go
type PredictiveScaler struct {
    predictor *LinearRegression  // 简单线性回归
}

func (s *PredictiveScaler) PredictLoad(history []MetricPoint, horizon time.Duration) (int, error) {
    // 基于历史负载预测
    if len(history) < 10 {
        return 0, errors.New("insufficient data")
    }

    // 训练模型
    s.predictor.Fit(history)

    // 预测未来负载
    futureLoad := s.predictor.Predict(time.Now().Add(horizon))

    // 计算所需实例数
    instances := int(math.Ceil(futureLoad / 100.0))  // 假设每个实例处理 100 RPS

    return instances, nil
}

func (s *PredictiveScaler) Scale(targetReplicas int) error {
    // 调用 K8s API 扩容
    return updateDeploymentReplicas(targetReplicas)
}
```

---

## 根因分析

```go
type RootCauseAnalyzer struct {
    graph DependencyGraph
}

func (a *RootCauseAnalyzer) Analyze(alert Alert) ([]Cause, error) {
    // 构建依赖图
    services := a.graph.GetDependencies(alert.Service)

    var causes []Cause

    // 检查上游服务
    for _, svc := range services {
        metrics := getServiceMetrics(svc)

        // 相关性分析
        correlation := calculateCorrelation(alert.Metrics, metrics)
        if correlation > 0.8 {
            causes = append(causes, Cause{
                Service:     svc,
                Correlation: correlation,
                TimeDelta:   alert.Time.Sub(getAnomalyTime(metrics)),
            })
        }
    }

    // 按时间排序，最早的异常可能是根因
    sort.Slice(causes, func(i, j int) bool {
        return causes[i].TimeDelta > causes[j].TimeDelta
    })

    return causes, nil
}
```
