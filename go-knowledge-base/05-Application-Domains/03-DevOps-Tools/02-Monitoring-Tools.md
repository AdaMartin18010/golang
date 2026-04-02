# 监控工具开发

> **分类**: 成熟应用领域

---

## Prometheus Exporter

```go
import "github.com/prometheus/client_golang/prometheus"

var (
    cpuTemp = prometheus.NewGauge(prometheus.GaugeOpts{
        Name: "cpu_temperature_celsius",
        Help: "Current temperature of the CPU.",
    })

    hdFailures = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "hd_errors_total",
            Help: "Number of hard-disk errors.",
        },
        []string{"device"},
    )
)

func init() {
    prometheus.MustRegister(cpuTemp, hdFailures)
}

func main() {
    go func() {
        for {
            cpuTemp.Set(getCPUTemperature())
            if isHDError() {
                hdFailures.WithLabelValues("/dev/sda").Inc()
            }
            time.Sleep(10 * time.Second)
        }
    }()

    http.Handle("/metrics", promhttp.Handler())
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

---

## 自定义 Collector

```go
type MyCollector struct {
    desc *prometheus.Desc
}

func NewMyCollector() *MyCollector {
    return &MyCollector{
        desc: prometheus.NewDesc("my_metric",
            "Description of my metric",
            nil, nil),
    }
}

func (c *MyCollector) Describe(ch chan<- *prometheus.Desc) {
    ch <- c.desc
}

func (c *MyCollector) Collect(ch chan<- prometheus.Metric) {
    value := getValue()
    ch <- prometheus.MustNewConstMetric(
        c.desc,
        prometheus.GaugeValue,
        value,
    )
}
```
