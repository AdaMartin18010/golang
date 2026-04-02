# Channel 高级模式

> **分类**: 开源技术堆栈  
> **标签**: #channel #advanced #patterns

---

## 有界并发

```go
func BoundedConcurrent(items []Item, maxConcurrent int) {
    sem := make(chan struct{}, maxConcurrent)
    var wg sync.WaitGroup
    
    for _, item := range items {
        wg.Add(1)
        go func(i Item) {
            defer wg.Done()
            
            sem <- struct{}{}        // 获取信号量
            defer func() { <-sem }()  // 释放
            
            process(i)
        }(item)
    }
    
    wg.Wait()
}
```

---

## 管道取消

```go
func PipelineWithCancel(ctx context.Context, stages ...Stage) Stage {
    return func(in <-chan int) <-chan int {
        // 在第一个阶段检查取消
        out := make(chan int)
        
        go func() {
            defer close(out)
            
            for {
                select {
                case <-ctx.Done():
                    // 排空输入
                    for range in {}
                    return
                case v, ok := <-in:
                    if !ok {
                        return
                    }
                    
                    select {
                    case out <- v:
                    case <-ctx.Done():
                        return
                    }
                }
            }
        }()
        
        return out
    }
}
```

---

## 扇出扇入控制

```go
// 可控扇出
func ControlledFanOut(ctx context.Context, source <-chan int, numWorkers int) []<-chan int {
    chans := make([]chan int, numWorkers)
    for i := range chans {
        chans[i] = make(chan int)
    }
    
    go func() {
        defer func() {
            for _, ch := range chans {
                close(ch)
            }
        }()
        
        i := 0
        for v := range source {
            select {
            case <-ctx.Done():
                return
            case chans[i%numWorkers] <- v:
                i++
            }
        }
    }()
    
    // 转换为只读
    result := make([]<-chan int, numWorkers)
    for i, ch := range chans {
        result[i] = ch
    }
    return result
}
```

---

## 或通道模式

```go
func Or(channels ...<-chan interface{}) <-chan interface{} {
    switch len(channels) {
    case 0:
        return nil
    case 1:
        return channels[0]
    }
    
    orDone := make(chan interface{})
    go func() {
        defer close(orDone)
        
        switch len(channels) {
        case 2:
            select {
            case <-channels[0]:
            case <-channels[1]:
            }
        default:
            select {
            case <-channels[0]:
            case <-channels[1]:
            case <-channels[2]:
            case <-Or(append(channels[3:], orDone)...):
            }
        }
    }()
    return orDone
}

// 使用：等待任意通道关闭
timeout := time.After(1 * time.Minute)
cancel := make(chan interface{})
<-Or(timeout, cancel)
```

---

## Tee 模式

```go
func Tee(ctx context.Context, in <-chan int) (<-chan int, <-chan int) {
    out1, out2 := make(chan int), make(chan int)
    
    go func() {
        defer close(out1)
        defer close(out2)
        
        for val := range in {
            var out1, out2 = out1, out2
            for i := 0; i < 2; i++ {
                select {
                case <-ctx.Done():
                    return
                case out1<-val:
                    out1 = nil
                case out2<-val:
                    out2 = nil
                }
            }
        }
    }()
    
    return out1, out2
}
```
