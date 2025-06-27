# HTTP协议基础

## 📚 **理论分析**

### **HTTP协议简介**

- HTTP（HyperText Transfer Protocol）是Web通信的基础协议，采用请求-响应模型。
- 无状态、基于文本、支持多种方法（GET、POST、PUT、DELETE等）。
- 主要用于客户端（浏览器/应用）与服务器之间的数据交换。

### **报文结构**

- **请求报文**：请求行、请求头、空行、请求体
- **响应报文**：状态行、响应头、空行、响应体

#### **请求报文示例**

```text
GET /index.html HTTP/1.1
Host: www.example.com
User-Agent: curl/7.68.0
Accept: */*

```

#### **响应报文示例**

```text
HTTP/1.1 200 OK
Content-Type: text/html; charset=UTF-8
Content-Length: 1024

<html>...</html>
```

### **常用HTTP方法**

- **GET**：获取资源
- **POST**：提交数据
- **PUT**：更新资源
- **DELETE**：删除资源
- **HEAD**：仅获取响应头
- **OPTIONS**：查询支持的方法

### **状态码分类**

- 1xx：信息（如100 Continue）
- 2xx：成功（如200 OK, 201 Created）
- 3xx：重定向（如301, 302）
- 4xx：客户端错误（如400, 404, 401）
- 5xx：服务器错误（如500, 502）

### **HTTP/1.1与HTTP/2对比**

- HTTP/1.1：串行请求、无多路复用、易受队头阻塞影响
- HTTP/2：多路复用、头部压缩、二进制分帧、性能更优

## 💻 **Go语言视角与代码示例**

### **发起HTTP请求（客户端）**

```go
package main
import (
    "fmt"
    "io/ioutil"
    "net/http"
)
func main() {
    resp, err := http.Get("https://httpbin.org/get")
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println(string(body))
}
```

### **解析HTTP请求（服务器）**

```go
package main
import (
    "fmt"
    "net/http"
)
func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Method: %s\nPath: %s\nUser-Agent: %s\n", r.Method, r.URL.Path, r.UserAgent())
}
func main() {
    http.HandleFunc("/", handler)
    http.ListenAndServe(":8080", nil)
}
```

## 🎯 **最佳实践**

- 使用HTTPS保障安全
- 合理设置超时与重试机制
- 正确处理状态码与错误
- 遵循RESTful API设计规范
- 日志记录请求与响应

## 🔍 **常见问题**

- Q: HTTP是有状态的吗？
  A: HTTP本身无状态，需用Cookie/Session等机制保持会话
- Q: 如何实现文件上传？
  A: 使用`multipart/form-data`编码，服务端解析
- Q: 如何防止XSS/CSRF？
  A: 输入校验、输出编码、CSRF Token

## 📚 **扩展阅读**

- [MDN HTTP协议详解](https://developer.mozilla.org/zh-CN/docs/Web/HTTP)
- [RFC 7230: HTTP/1.1 Message Syntax](https://tools.ietf.org/html/rfc7230)
- [Go net/http官方文档](https://golang.org/pkg/net/http/)

---

**文档维护者**: AI Assistant  
**最后更新**: 2024年6月27日  
**文档状态**: 完成
