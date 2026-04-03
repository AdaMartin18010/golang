# AD-031-IoT-Edge-Computing-2026

> **Dimension**: 05-Application-Domains  
> **Status**: S-Level  
> **Created**: 2026-04-03  
> **Version**: 2026 (MQTT, CoAP, Edge ML, Time-Series DB)  
> **Size**: >20KB 

---

## 1. IoT架构概览

### 1.1 分层架构

```
┌─────────────────────────────────────────┐
│          Cloud Layer                    │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐   │
│  │ Analytics│ │ ML     │ │ Storage │   │
│  └─────────┘ └─────────┘ └─────────┘   │
└─────────────────────────────────────────┘
                    │
┌─────────────────────────────────────────┐
│          Fog Layer                      │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐   │
│  │ Gateway │ │ Stream  │ │ Local   │   │
│  │         │ │ Process │ │ Storage │   │
│  └─────────┘ └─────────┘ └─────────┘   │
└─────────────────────────────────────────┘
                    │
┌─────────────────────────────────────────┐
│          Edge Layer                     │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐   │
│  │ Sensors │ │ Actuators│ │ Edge    │   │
│  │ Devices │ │         │ │ Compute │   │
│  └─────────┘ └─────────┘ └─────────┘   │
└─────────────────────────────────────────┘
```

### 1.2 协议对比

| 协议 | 传输层 | 功耗 | 适用场景 |
|------|--------|------|---------|
| MQTT | TCP | 中 | 可靠传输 |
| CoAP | UDP | 低 | 受限设备 |
| HTTP/2 | TCP | 高 | 高带宽 |
| LoRaWAN | RF | 极低 | 远距离低功耗 |
| Zigbee | RF | 低 | 短距离Mesh |

---

## 2. MQTT协议

### 2.1 协议特性

```
MQTT架构:

Publisher ──► Broker ◄── Subscriber
   │              │
   └──────────────┘
      Topic: sensors/temperature/room1

QoS级别:
0 - 最多一次 (fire and forget)
1 - 至少一次 (acknowledged)
2 - 恰好一次 (assured delivery)

Topic通配符:
+ - 单级通配符 (sensors/+/temperature)
# - 多级通配符 (sensors/#)
```

### 2.2 Go实现

```go
// MQTT Publisher
package main

import (
    mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
    opts := mqtt.NewClientOptions()
    opts.AddBroker("tcp://broker.hivemq.com:1883")
    opts.SetClientID("iot-sensor-001")
    opts.SetUsername("user")
    opts.SetPassword("pass")
    
    // 连接回调
    opts.OnConnect = func(c mqtt.Client) {
        log.Println("Connected to MQTT broker")
    }
    
    // 断线重连
    opts.SetAutoReconnect(true)
    opts.SetConnectRetry(true)
    
    client := mqtt.NewClient(opts)
    if token := client.Connect(); token.Wait() && token.Error() != nil {
        log.Fatal(token.Error())
    }
    
    // 发布消息
    sensorData := map[string]interface{}{
        "device_id":   "sensor-001",
        "temperature": 23.5,
        "humidity":    65,
        "timestamp":   time.Now().Unix(),
    }
    
    payload, _ := json.Marshal(sensorData)
    
    token := client.Publish(
        "factory/zone-a/temperature",
        1,      // QoS
        false,  // Retain
        payload,
    )
    token.Wait()
    
    client.Disconnect(250)
}

// MQTT Subscriber
func subscriber() {
    opts := mqtt.NewClientOptions()
    opts.AddBroker("tcp://broker.hivemq.com:1883")
    opts.SetClientID("iot-consumer-001")
    
    // 消息回调
    opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
        log.Printf("Topic: %s, Message: %s\n", msg.Topic(), msg.Payload())
        
        // 解析处理
        var data SensorData
        if err := json.Unmarshal(msg.Payload(), &data); err != nil {
            log.Printf("Parse error: %v", err)
            return
        }
        
        // 存储到时序数据库
        storeToTimeseriesDB(data)
    })
    
    client := mqtt.NewClient(opts)
    if token := client.Connect(); token.Wait() && token.Error() != nil {
        log.Fatal(token.Error())
    }
    
    // 订阅多个主题
    topics := map[string]byte{
        "factory/+/temperature": 1,  // QoS 1
        "factory/+/humidity":    1,
        "factory/alerts/#":      2,  // QoS 2 for alerts
    }
    
    if token := client.SubscribeMultiple(topics, nil); token.Wait() && token.Error() != nil {
        log.Fatal(token.Error())
    }
    
    // 保持运行
    select {}
}
```

### 2.3 MQTT 5.0新特性

```go
// MQTT 5.0 增强
opts.SetProtocolVersion(5)

// 共享订阅 (负载均衡)
token := client.Subscribe("$share/group1/factory/+/temperature", 1, handler)

// 消息过期
msg := &mqtt.PublishOptions{
    Payload: payload,
    Properties: &mqtt.PublishProperties{
        MessageExpiry: 60,  // 60秒后过期
    },
}

// 流量控制
opts.SetReceiveMaximum(100)  // 同时处理100条消息

// 用户属性
props := &mqtt.UserProperties{}
props.Add("sensor-type", "temperature")
props.Add("zone", "A")
```

---

## 3. CoAP协议

### 3.1 协议特性

```
CoAP (Constrained Application Protocol):
- 基于UDP
- 低功耗
- RESTful API
- 支持观察模式

消息类型:
CON - Confirmable (需要确认)
NON - Non-confirmable (不需要确认)
ACK - Acknowledgment
RST - Reset
```

### 3.2 Go实现

```go
package main

import (
    "github.com/plgd-dev/go-coap/v3"
    "github.com/plgd-dev/go-coap/v3/message"
)

// CoAP Server
func startCoAPServer() {
    mux := coap.NewServeMux()
    
    mux.HandleFunc("/sensors/temperature", func(w coap.ResponseWriter, r *coap.Request) {
        switch r.Code {
        case codes.GET:
            // 读取传感器数据
            temp := readTemperature()
            w.SetContentFormat(message.AppJSON)
            w.Write([]byte(fmt.Sprintf(`{"temperature":%.2f}`, temp)))
            
        case codes.POST:
            // 配置传感器
            config := parseConfig(r.Body())
            configureSensor(config)
            w.SetCode(codes.Changed)
        }
    })
    
    // 观察模式 (类似MQTT订阅)
    mux.HandleFunc("/sensors/temperature/observe", func(w coap.ResponseWriter, r *coap.Request) {
        obs, err := w.Observe()
        if err != nil {
            return
        }
        
        ticker := time.NewTicker(5 * time.Second)
        defer ticker.Stop()
        
        for range ticker.C {
            if !obs.IsActive() {
                return
            }
            
            temp := readTemperature()
            obs.Notify([]byte(fmt.Sprintf(`{"temperature":%.2f}`, temp)))
        }
    })
    
    log.Fatal(coap.ListenAndServe("udp", ":5683", mux))
}

// CoAP Client
func coapClient() {
    conn, err := coap.Dial("udp", "coap-server:5683")
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()
    
    // GET请求
    resp, err := conn.Get(ctx, "/sensors/temperature")
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Temperature: %s", resp.Body())
    
    // 观察模式
    obs, err := conn.Observe(ctx, "/sensors/temperature/observe", func(req *coap.Request) {
        log.Printf("Update: %s", req.Body())
    })
    if err != nil {
        log.Fatal(err)
    }
    
    time.Sleep(1 * time.Minute)
    obs.Cancel()
}
```

---

## 4. 边缘计算

### 4.1 边缘架构

```go
// Edge Gateway
type EdgeGateway struct {
    deviceManager  *DeviceManager
    ruleEngine     *RuleEngine
    localStorage   *LocalStorage
    cloudConnector *CloudConnector
    mlRuntime      *MLRuntime
}

func (g *EdgeGateway) Start() {
    // 设备管理
    go g.deviceManager.Run()
    
    // 规则引擎
    go g.ruleEngine.Run()
    
    // 数据转发
    go g.dataForwarder()
    
    // ML推理服务
    go g.mlRuntime.Serve()
}

func (g *EdgeGateway) dataForwarder() {
    batch := make([]SensorData, 0, 100)
    ticker := time.NewTicker(30 * time.Second)
    
    for {
        select {
        case data := <-g.deviceManager.DataChan:
            // 本地预处理
            processed := g.preprocess(data)
            
            // 检查规则
            if g.ruleEngine.Evaluate(processed) {
                g.executeAction(processed)
            }
            
            batch = append(batch, processed)
            
            // 本地存储
            g.localStorage.Write(processed)
            
        case <-ticker.C:
            if len(batch) > 0 {
                // 批量上传到云端
                g.cloudConnector.UploadBatch(batch)
                batch = batch[:0]
            }
        }
    }
}
```

### 4.2 边缘ML推理

```go
// TensorFlow Lite推理
type EdgeML struct {
    interpreter *tflite.Interpreter
}

func NewEdgeML(modelPath string) (*EdgeML, error) {
    model, err := tflite.LoadModel(modelPath)
    if err != nil {
        return nil, err
    }
    
    options := tflite.NewInterpreterOptions()
    options.SetNumThreads(4)
    options.SetDelegate(tflite.GPUDelegate)  // 使用GPU加速
    
    interpreter, err := tflite.NewInterpreter(model, options)
    if err != nil {
        return nil, err
    }
    
    return &EdgeML{interpreter: interpreter}, nil
}

func (m *EdgeML) Predict(sensorData []float32) (Prediction, error) {
    input := m.interpreter.GetInputTensor(0)
    input.CopyFromBuffer(sensorData)
    
    if err := m.interpreter.Invoke(); err != nil {
        return Prediction{}, err
    }
    
    output := m.interpreter.GetOutputTensor(0)
    var result []float32
    output.CopyToBuffer(&result)
    
    return Prediction{
        Anomaly: result[0] > 0.8,
        Score:   result[0],
    }, nil
}

// 异常检测流程
func (g *EdgeGateway) anomalyDetection(data SensorData) {
    features := extractFeatures(data)
    
    prediction, err := g.mlRuntime.Predict(features)
    if err != nil {
        return
    }
    
    if prediction.Anomaly {
        // 本地告警
        g.triggerLocalAlert(data, prediction)
        
        // 上报云端
        g.cloudConnector.ReportAnomaly(data, prediction)
    }
}
```

---

## 5. 时序数据库

### 5.1 InfluxDB集成

```go
package main

import (
    "github.com/influxdata/influxdb-client-go/v2"
    "github.com/influxdata/influxdb-client-go/v2/api/write"
)

func main() {
    client := influxdb2.NewClient("http://localhost:8086", "token")
    defer client.Close()
    
    writeAPI := client.WriteAPI("org", "bucket")
    
    // 异步写入
    writeAPI.WritePoint(
        write.NewPoint(
            "temperature",
            map[string]string{
                "device": "sensor-001",
                "zone":   "A",
            },
            map[string]interface{}{
                "value": 23.5,
                "unit":  "celsius",
            },
            time.Now(),
        ),
    )
    
    // 刷新缓冲区
    writeAPI.Flush()
    
    // 查询数据
    queryAPI := client.QueryAPI("org")
    result, err := queryAPI.Query(context.Background(), `
        from(bucket: "bucket")
            |> range(start: -1h)
            |> filter(fn: (r) => r._measurement == "temperature")
            |> aggregateWindow(every: 5m, fn: mean)
    `)
    
    if err != nil {
        log.Fatal(err)
    }
    
    for result.Next() {
        log.Printf("Time: %v, Value: %v", result.Record().Time(), result.Record().Value())
    }
}
```

### 5.2 数据降采样

```go
// 边缘侧数据降采样
func (g *EdgeGateway) downsample(data []SensorData) []AggregatedData {
    buckets := make(map[int64][]SensorData)
    
    // 按5分钟桶分组
    for _, d := range data {
        bucket := d.Timestamp.Unix() / 300 * 300
        buckets[bucket] = append(buckets[bucket], d)
    }
    
    result := make([]AggregatedData, 0, len(buckets))
    
    for timestamp, group := range buckets {
        var sum float64
        var min, max float64 = group[0].Value, group[0].Value
        
        for _, d := range group {
            sum += d.Value
            if d.Value < min {
                min = d.Value
            }
            if d.Value > max {
                max = d.Value
            }
        }
        
        result = append(result, AggregatedData{
            Timestamp: time.Unix(timestamp, 0),
            Avg:       sum / float64(len(group)),
            Min:       min,
            Max:       max,
            Count:     len(group),
        })
    }
    
    return result
}
```

---

## 6. 安全

### 6.1 设备认证

```go
// X.509证书认证
type DeviceAuth struct {
    caCert     *x509.Certificate
    deviceCert *x509.Certificate
    privateKey *ecdsa.PrivateKey
}

func (a *DeviceAuth) MutualTLSConfig() *tls.Config {
    cert := tls.Certificate{
        Certificate: [][]byte{a.deviceCert.Raw},
        PrivateKey:  a.privateKey,
    }
    
    caPool := x509.NewCertPool()
    caPool.AddCert(a.caCert)
    
    return &tls.Config{
        Certificates:       []tls.Certificate{cert},
        RootCAs:           caPool,
        InsecureSkipVerify: false,
        VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
            // 额外验证设备ID
            return verifyDeviceCertificate(rawCerts)
        },
    }
}

// JWT设备认证
func (g *EdgeGateway) authenticateDevice(deviceID string, token string) error {
    parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
        return g.publicKey, nil
    })
    
    if err != nil || !parsedToken.Valid {
        return ErrInvalidToken
    }
    
    claims := parsedToken.Claims.(jwt.MapClaims)
    if claims["device_id"] != deviceID {
        return ErrDeviceMismatch
    }
    
    return nil
}
```

### 6.2 数据加密

```go
// 敏感数据加密传输
func encryptPayload(data []byte, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }
    
    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return nil, err
    }
    
    return gcm.Seal(nonce, nonce, data, nil), nil
}
```

---

## 7. 最佳实践

### 7.1 网络优化

| 优化点 | 方案 |
|--------|------|
| 带宽限制 | 数据压缩、差分传输 |
| 网络不稳定 | 本地缓存、断点续传 |
| 高延迟 | 边缘预处理、批量传输 |
| 功耗限制 | 休眠调度、数据聚合 |

### 7.2 设备管理

```go
// OTA更新
type OTAManager struct {
    currentVersion string
    updateServer   string
}

func (o *OTAManager) CheckUpdate() (*FirmwareUpdate, error) {
    resp, err := http.Get(fmt.Sprintf("%s/firmware/latest?current=%s", 
        o.updateServer, o.currentVersion))
    if err != nil {
        return nil, err
    }
    
    var update FirmwareUpdate
    if err := json.NewDecoder(resp.Body).Decode(&update); err != nil {
        return nil, err
    }
    
    return &update, nil
}

func (o *OTAManager) ApplyUpdate(update *FirmwareUpdate) error {
    // 下载固件
    firmware, err := o.downloadFirmware(update.URL)
    if err != nil {
        return err
    }
    
    // 验证签名
    if !o.verifySignature(firmware, update.Signature) {
        return ErrInvalidSignature
    }
    
    // 备份当前固件
    o.backupCurrent()
    
    // 写入新固件
    if err := o.writeFirmware(firmware); err != nil {
        o.rollback()
        return err
    }
    
    // 重启
    o.reboot()
    
    return nil
}
```

---

## 8. 参考文献

1. MQTT Specification 5.0
2. CoAP RFC 7252
3. "Edge Computing" - Weisong Shi
4. InfluxDB Documentation
5. TensorFlow Lite Guide

---

*Last Updated: 2026-04-03*
