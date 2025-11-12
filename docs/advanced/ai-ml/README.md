# Go AI/ML开发

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---
## 📋 目录

- [Go AI/ML开发](#go-aiml开发)
  - [📚 核心内容](#核心内容)
  - [🚀 TensorFlow Serving](#tensorflow-serving)
  - [📖 系统文档](#系统文档)

---

## 📚 核心内容

1. **[Go与AI集成](./01-Go与AI集成.md)** ⭐⭐⭐⭐
2. **[机器学习库](./02-机器学习库.md)** ⭐⭐⭐⭐
3. **[深度学习框架](./03-深度学习框架.md)** ⭐⭐⭐
4. **[模型推理](./04-模型推理.md)** ⭐⭐⭐⭐⭐
5. **[数据处理](./05-数据处理.md)** ⭐⭐⭐⭐
6. **[实战案例](./06-实战案例.md)** ⭐⭐⭐⭐

---

## 🚀 TensorFlow Serving

```go
import tf "github.com/tensorflow/tensorflow/tensorflow/go"

// 加载模型
model, _ := tf.LoadSavedModel("model_path", []string{"serve"}, nil)

// 推理
result, _ := model.Session.Run(
    map[tf.Output]*tf.Tensor{
        model.Graph.Operation("input").Output(0): inputTensor,
    },
    []tf.Output{
        model.Graph.Operation("output").Output(0),
    },
    nil,
)
```

---

## 📖 系统文档
