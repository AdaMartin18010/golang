# Go与AI集成

**版本**: v1.0  
**更新日期**: 2025-10-29  
**适用于**: Go 1.25.3

---

## 📋 目录



- [1. 📖 概念介绍](#1-概念介绍)
- [2. 🎯 为什么使用Go做AI](#2-为什么使用go做ai)
  - [优势](#优势)
  - [适用场景](#适用场景)
- [🔧 调用Python模型](#调用python模型)
  - [1. 通过HTTP API](#1-通过http-api)
  - [2. 通过gRPC](#2-通过grpc)
- [🚀 TensorFlow集成](#tensorflow集成)
- [🔮 ONNX Runtime](#onnx-runtime)
- [🌐 构建推理服务](#构建推理服务)
- [💡 最佳实践](#最佳实践)
  - [1. 批量推理](#1-批量推理)
  - [2. 模型缓存](#2-模型缓存)
- [📚 相关资源](#相关资源)

## 1. 📖 概念介绍

Go在AI领域的应用日益增长，特别是在模型部署、推理服务和数据处理方面。本文介绍如何将Go与AI/ML生态系统集成。

---

## 2. 🎯 为什么使用Go做AI

### 优势
- **高性能**: 编译型语言，接近C的性能
- **并发**: 天然支持高并发，适合推理服务
- **部署简单**: 单一二进制文件，无依赖
- **云原生**: 容器化部署的理想选择

### 适用场景
- 模型推理服务
- 数据预处理管道
- AI模型API网关
- 实时预测系统

---

## 🔧 调用Python模型

### 1. 通过HTTP API

```go
package main

import (
    "bytes"
    "encoding/json"
    "net/http"
)

type PredictionRequest struct {
    Data []float64 `json:"data"`
}

type PredictionResponse struct {
    Prediction float64 `json:"prediction"`
    Confidence float64 `json:"confidence"`
}

// 调用Python Flask模型服务
func callPythonModel(data []float64) (*PredictionResponse, error) {
    url := "http://localhost:5000/predict"
    
    reqBody := PredictionRequest{Data: data}
    jsonData, _ := json.Marshal(reqBody)
    
    resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var result PredictionResponse
    json.NewDecoder(resp.Body).Decode(&result)
    
    return &result, nil
}

// 使用示例
func main() {
    data := []float64{1.0, 2.0, 3.0, 4.0}
    result, err := callPythonModel(data)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Prediction: %.2f (Confidence: %.2f%%)\n", 
        result.Prediction, result.Confidence*100)
}
```

---

### 2. 通过gRPC

```protobuf
// prediction.proto
syntax = "proto3";

service Predictor {
  rpc Predict(PredictRequest) returns (PredictResponse);
}

message PredictRequest {
  repeated float features = 1;
}

message PredictResponse {
  float prediction = 1;
  float confidence = 2;
}
```

```go
// Go客户端
import (
    "context"
    pb "path/to/prediction"
    "google.golang.org/grpc"
)

func callGRPCModel(features []float32) (*pb.PredictResponse, error) {
    conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
    if err != nil {
        return nil, err
    }
    defer conn.Close()
    
    client := pb.NewPredictorClient(conn)
    
    req := &pb.PredictRequest{Features: features}
    resp, err := client.Predict(context.Background(), req)
    
    return resp, err
}
```

---

## 🚀 TensorFlow集成

```go
import (
    tf "github.com/tensorflow/tensorflow/tensorflow/go"
    "github.com/tensorflow/tensorflow/tensorflow/go/op"
)

type TFModel struct {
    model   *tf.SavedModel
    session *tf.Session
}

func LoadTFModel(modelPath string) (*TFModel, error) {
    model, err := tf.LoadSavedModel(modelPath, []string{"serve"}, nil)
    if err != nil {
        return nil, err
    }
    
    return &TFModel{
        model:   model,
        session: model.Session,
    }, nil
}

func (m *TFModel) Predict(input [][]float32) ([]float32, error) {
    // 创建输入tensor
    tensor, err := tf.NewTensor(input)
    if err != nil {
        return nil, err
    }
    
    // 运行推理
    results, err := m.session.Run(
        map[tf.Output]*tf.Tensor{
            m.model.Graph.Operation("input").Output(0): tensor,
        },
        []tf.Output{
            m.model.Graph.Operation("output").Output(0),
        },
        nil,
    )
    
    if err != nil {
        return nil, err
    }
    
    return results[0].Value().([]float32), nil
}

// 使用示例
func example() {
    model, _ := LoadTFModel("./saved_model")
    defer model.session.Close()
    
    input := [][]float32{{1.0, 2.0, 3.0}}
    output, _ := model.Predict(input)
    
    fmt.Printf("Prediction: %v\n", output)
}
```

---

## 🔮 ONNX Runtime

```go
import "github.com/yalue/onnxruntime_go"

type ONNXModel struct {
    session *onnxruntime.Session
    input   *onnxruntime.Tensor
    output  *onnxruntime.Tensor
}

func LoadONNXModel(modelPath string) (*ONNXModel, error) {
    // 初始化ONNX Runtime
    err := onnxruntime.InitializeEnvironment()
    if err != nil {
        return nil, err
    }
    
    // 加载模型
    session, err := onnxruntime.NewSession(modelPath, nil)
    if err != nil {
        return nil, err
    }
    
    return &ONNXModel{session: session}, nil
}

func (m *ONNXModel) Predict(data []float32) ([]float32, error) {
    // 创建输入tensor
    inputShape := []int64{1, int64(len(data))}
    inputTensor, err := onnxruntime.NewTensor(inputShape, data)
    if err != nil {
        return nil, err
    }
    defer inputTensor.Destroy()
    
    // 运行推理
    outputs, err := m.session.Run([]onnxruntime.Value{inputTensor})
    if err != nil {
        return nil, err
    }
    defer outputs[0].Destroy()
    
    // 提取结果
    outputTensor := outputs[0].GetTensor()
    return outputTensor.GetData().([]float32), nil
}
```

---

## 🌐 构建推理服务

```go
package main

import (
    "encoding/json"
    "net/http"
    "sync"
)

type ModelServer struct {
    model *ONNXModel
    pool  *sync.Pool
}

func NewModelServer(modelPath string) (*ModelServer, error) {
    model, err := LoadONNXModel(modelPath)
    if err != nil {
        return nil, err
    }
    
    return &ModelServer{
        model: model,
        pool: &sync.Pool{
            New: func() interface{} {
                return make([]float32, 0, 100)
            },
        },
    }, nil
}

func (ms *ModelServer) PredictHandler(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Features []float32 `json:"features"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // 推理
    result, err := ms.model.Predict(req.Features)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    // 返回结果
    json.NewEncoder(w).Encode(map[string]interface{}{
        "prediction": result,
    })
}

func main() {
    server, _ := NewModelServer("model.onnx")
    
    http.HandleFunc("/predict", server.PredictHandler)
    http.ListenAndServe(":8080", nil)
}
```

---

## 💡 最佳实践

### 1. 批量推理
```go
type BatchProcessor struct {
    batchSize int
    timeout   time.Duration
    queue     chan Request
}

func (bp *BatchProcessor) Process() {
    batch := make([]Request, 0, bp.batchSize)
    timer := time.NewTimer(bp.timeout)
    
    for {
        select {
        case req := <-bp.queue:
            batch = append(batch, req)
            if len(batch) >= bp.batchSize {
                bp.processBatch(batch)
                batch = batch[:0]
                timer.Reset(bp.timeout)
            }
        case <-timer.C:
            if len(batch) > 0 {
                bp.processBatch(batch)
                batch = batch[:0]
            }
            timer.Reset(bp.timeout)
        }
    }
}
```

### 2. 模型缓存
```go
type ModelCache struct {
    models map[string]*ONNXModel
    mu     sync.RWMutex
}

func (mc *ModelCache) Get(modelID string) (*ONNXModel, error) {
    mc.mu.RLock()
    model, exists := mc.models[modelID]
    mc.mu.RUnlock()
    
    if exists {
        return model, nil
    }
    
    // 加载模型
    mc.mu.Lock()
    defer mc.mu.Unlock()
    
    model, err := LoadONNXModel(modelID)
    if err != nil {
        return nil, err
    }
    
    mc.models[modelID] = model
    return model, nil
}
```

---

## 📚 相关资源

- [TensorFlow Go](https://github.com/tensorflow/tensorflow/tree/master/tensorflow/go)
- [ONNX Runtime Go](https://github.com/yalue/onnxruntime_go)
- [GoLearn](https://github.com/sjwhitworth/golearn)

**下一步**: [02-机器学习库](./02-机器学习库.md)

---

**最后更新**: 2025-10-29

