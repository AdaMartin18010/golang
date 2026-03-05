# Go错误处理哲学与实践

> 从软件工程角度理解Go的错误处理设计

---

## 一、错误处理的工程考量

### 1.1 异常 vs 返回值的本质区别

```text
控制流的可见性：
────────────────────────────────────────

异常的问题：

// Java代码
public void process() {
    doSomething();  // 可能抛出异常？不知道
    doAnother();    // 可能抛出异常？不知道
    doFinal();      // 可能抛出异常？不知道
}

异常可以：
1. 在任何地方抛出
2. 在调用栈的任何层级捕获
3. 完全不可见地从函数中"跳出"

结果是：
- 代码执行路径难以预测
- 资源清理困难（需要try-finally）
- 错误处理分散在各处

Go的显式返回：
────────────────────────────────────────

func process() error {
    if err := doSomething(); err != nil {
        return err  // 错误明确可见
    }
    if err := doAnother(); err != nil {
        return err
    }
    if err := doFinal(); err != nil {
        return err
    }
    return nil
}

优势：
- 错误路径清晰可见
- 控制流线性，易于理解
- 资源清理简单（defer）

实际案例：
────────────────────────────────────────

某Java系统的问题：
- 一个NullPointerException在底层抛出
- 在中间层被捕获并包装
- 在顶层再次捕获并记录
- 真正的错误原因被多层包装掩盖
- 花了数小时才找到根本原因

Go版本：
- 每层返回错误时添加上下文
- 错误链清晰完整
- 一眼就能看出问题所在
```

### 1.2 错误处理的性能考量

```text
异常的性能陷阱：
────────────────────────────────────────

创建异常对象昂贵：
- 需要捕获完整的堆栈跟踪
- 涉及多次内存分配
- 在热点代码中影响显著

Java示例：
// 在循环中抛出异常
for (int i = 0; i < 100000; i++) {
    try {
        parseNumber(str);  // 可能抛出ParseException
    } catch (ParseException e) {
        // 处理
    }
}
// 每次异常创建成本约 1-5 μs

Go的错误值：
────────────────────────────────────────

// 简单的值返回，无额外开销
for i := 0; i < 100000; i++ {
    if _, err := parseNumber(str); err != nil {
        // 处理
    }
}
// 错误检查成本约 1-2 ns

性能对比：
────────────────────────────────────────

操作                    耗时
────────────────────────────────
Go错误返回              1-2 ns
Java异常（不创建栈）     50-100 ns
Java异常（创建栈）       1-5 μs
C++异常                 1-10 μs

在错误频繁的路径上：
- 异常机制可能消耗大量CPU
- 错误值几乎无开销

Go的设计哲学：
"错误是常态，不是异常"
────────────────────────────────────────

在Go中：
- 文件不存在是错误，不是异常
- 网络超时是错误，不是异常
- 无效输入是错误，不是异常

这些情况下使用返回值是合理的。

异常适合的真正场景：
- 编程错误（如数组越界）
- 不可恢复的系统错误
- Go使用panic处理这些情况
```

---

## 二、错误包装的艺术

### 2.1 为什么需要错误链

```text
错误上下文的丢失：
────────────────────────────────────────

场景：HTTP请求处理失败

不包装的错误：
"connection refused"

问题：
- 哪个服务连接被拒绝？
- 在处理哪个请求时发生？
- 是数据库？缓存？还是第三方API？

包装后的错误：
"handle /api/users: load user 123: query database: connection refused"

现在我们知道：
- 在处理 /api/users 请求
- 加载ID为123的用户时
- 查询数据库时
- 连接被拒绝

错误包装的价值：
────────────────────────────────────────

1. 定位问题：
   快速确定错误发生的具体位置

2. 理解上下文：
   了解导致错误的完整调用链

3. 调试友好：
   日志中的错误信息完整清晰

4. 分类处理：
   可以根据不同层级的错误做不同处理
```

### 2.2 错误包装的最佳实践

```text
包装原则：
────────────────────────────────────────

1. 添加上下文，不要隐藏原因：

// 不良：隐藏了底层错误
if err != nil {
    return errors.New("database error")
}

// 良好：保留原始错误
if err != nil {
    return fmt.Errorf("query user %d: %w", userID, err)
}

2. 使用 %w 保留错误链：

// 可以被 errors.Is 和 errors.As 识别
return fmt.Errorf("process order: %w", err)

3. 在关键边界包装：

- 服务边界（API、数据库）
- 层边界（handler → service → repository）
- 重大操作边界

不需要在每个函数都包装：
────────────────────────────────────────

// 过度包装
func a() error {
    if err := b(); err != nil {
        return fmt.Errorf("a failed: %w", err)
    }
    return nil
}

func b() error {
    if err := c(); err != nil {
        return fmt.Errorf("b failed: %w", err)
    }
    return nil
}

// 导致错误链过长：
// "a failed: b failed: c failed: connection refused"

// 更合理的方式：只在有意义的边界包装
func serviceLayer() error {
    if err := repo.Query(); err != nil {
        return fmt.Errorf("load user: %w", err)
    }
    return nil
}

错误消息的风格：
────────────────────────────────────────

1. 简洁明了：
   "load user: %w" 而不是 "failed to load user with error: %w"

2. 使用小写开头：
   错误消息通常被包裹在其他句子中

3. 包含关键信息：
   用户ID、资源名称等有助于定位的信息

4. 避免冗余：
   不需要 "error occurred while..."
```

---

## 三、Sentinel错误与错误类型

### 3.1 预定义错误 (Sentinel Errors)

```text
什么是Sentinel错误：
────────────────────────────────────────

预定义的错误值，用于特定的错误情况。

标准库示例：
var (
    ErrNotExist   = errors.New("file does not exist")
    ErrPermission = errors.New("permission denied")
    ErrExist      = errors.New("file already exists")
)

使用场景：
────────────────────────────────────────

场景1：区分不同错误情况

func processFile(path string) error {
    f, err := os.Open(path)
    if err != nil {
        if errors.Is(err, os.ErrNotExist) {
            // 文件不存在，创建它
            return createFile(path)
        }
        return err
    }
    defer f.Close()
    // 处理文件
}

场景2：重试逻辑

func callAPI() error {
    for i := 0; i < 3; i++ {
        err := apiRequest()
        if err == nil {
            return nil
        }
        if errors.Is(err, ErrRateLimited) {
            // 限流，等待后重试
            time.Sleep(time.Second * time.Duration(i+1))
            continue
        }
        // 其他错误，不重试
        return err
    }
    return ErrMaxRetries
}

定义Sentinel错误：
────────────────────────────────────────

package mypkg

var (
    ErrUserNotFound    = errors.New("user not found")
    ErrInvalidInput    = errors.New("invalid input")
    ErrUnauthorized    = errors.New("unauthorized")
    ErrInternal        = errors.New("internal error")
)

// 在错误链中识别
if errors.Is(err, ErrUserNotFound) {
    http.Error(w, "User not found", http.StatusNotFound)
}
```

### 3.2 自定义错误类型

```text
什么时候需要自定义类型：
────────────────────────────────────────

当需要额外的错误信息时：

// Sentinel错误无法携带详细信息
// 只知道是"validation error"，不知道哪个字段

自定义错误类型：

type ValidationError struct {
    Field   string
    Message string
    Value   interface{}
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error: field %s: %s", e.Field, e.Message)
}

使用：
────────────────────────────────────────

func validateUser(user *User) error {
    if user.Email == "" {
        return &ValidationError{
            Field:   "email",
            Message: "required",
            Value:   user.Email,
        }
    }
    // ...
}

处理：

func handleRequest(w http.ResponseWriter, r *http.Request) {
    if err := process(r); err != nil {
        var valErr *ValidationError
        if errors.As(err, &valErr) {
            // 知道具体是哪个字段错误
            json.NewEncoder(w).Encode(map[string]string{
                "error": "validation failed",
                "field": valErr.Field,
            })
            return
        }
        http.Error(w, err.Error(), 500)
    }
}

更复杂的错误类型：
────────────────────────────────────────

// 包含多个验证错误
type MultiValidationError struct {
    Errors []ValidationError
}

func (e *MultiValidationError) Error() string {
    var msgs []string
    for _, err := range e.Errors {
        msgs = append(msgs, err.Error())
    }
    return strings.Join(msgs, "; ")
}

func (e *MultiValidationError) Add(err ValidationError) {
    e.Errors = append(e.Errors, err)
}

func (e *MultiValidationError) HasErrors() bool {
    return len(e.Errors) > 0
}
```

---

## 四、错误处理模式

### 4.1 优雅的错误处理

```text
模式1：尽早返回
────────────────────────────────────────

// 不良：嵌套层级深
func process(data *Data) error {
    if data != nil {
        if data.Field != "" {
            if valid(data.Field) {
                // 实际处理
                return nil
            } else {
                return errors.New("invalid field")
            }
        } else {
            return errors.New("empty field")
        }
    } else {
        return errors.New("nil data")
    }
}

// 良好：扁平化
func process(data *Data) error {
    if data == nil {
        return errors.New("nil data")
    }
    if data.Field == "" {
        return errors.New("empty field")
    }
    if !valid(data.Field) {
        return errors.New("invalid field")
    }
    // 实际处理
    return nil
}

模式2：包装与转发
────────────────────────────────────────

func serviceLayer() error {
    if err := repo.Query(); err != nil {
        return fmt.Errorf("database operation failed: %w", err)
    }
    return nil
}

func handlerLayer() error {
    if err := serviceLayer(); err != nil {
        // 添加更上层的上下文
        return fmt.Errorf("handle request: %w", err)
    }
    return nil
}

模式3：使用辅助函数
────────────────────────────────────────

// 检查并包装
func wrap(err error, msg string) error {
    if err != nil {
        return fmt.Errorf("%s: %w", msg, err)
    }
    return nil
}

// 使用
func process() error {
    if err := step1(); err != nil {
        return wrap(err, "step1")
    }
    if err := step2(); err != nil {
        return wrap(err, "step2")
    }
    return nil
}
```

---

*本章深入探讨了Go错误处理的设计哲学和最佳实践。*
