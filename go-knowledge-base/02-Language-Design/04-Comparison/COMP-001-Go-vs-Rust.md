# Go vs Rust: Comprehensive Language Comparison

## Executive Summary

Go and Rust represent two distinct philosophies in systems programming: Go prioritizes simplicity and developer productivity with a garbage collector, while Rust emphasizes memory safety without garbage collection through its ownership system. This document provides an in-depth comparison across multiple dimensions.

---

## Table of Contents

1. [Language Philosophy](#language-philosophy)
2. [Concurrency Models](#concurrency-models)
3. [Memory Management](#memory-management)
4. [Performance Characteristics](#performance-characteristics)
5. [Code Examples](#code-examples)
6. [Performance Benchmarks](#performance-benchmarks)
7. [Decision Matrix](#decision-matrix)
8. [Migration Guide](#migration-guide)
9. [When to Choose Which](#when-to-choose-which)

---

## Language Philosophy

### Go Philosophy

```
Simplicity is the ultimate sophistication.
                                    — Leonardo da Vinci
```

Go was designed at Google with explicit goals:
- **Simplicity**: Minimal syntax, orthogonal features
- **Readability**: Code should be obvious to any Go programmer
- **Fast compilation**: Scale to millions of lines quickly
- **Efficient execution**: Comparable to C++ for many tasks
- **Great concurrency**: First-class support for concurrent programming

### Rust Philosophy

```
With great power comes great responsibility.
                                    — Spider-Man (and Rust)
```

Rust was designed by Mozilla with core principles:
- **Memory safety without GC**: Zero-cost abstractions
- **Fearless concurrency**: Compile-time data race prevention
- **Zero-cost abstractions**: High-level features without runtime overhead
- **Systems programming**: Suitable for kernels, embedded, game engines

---

## Concurrency Models

### Go: CSP (Communicating Sequential Processes)

Go implements concurrency through goroutines and channels based on Hoare's CSP model.

#### Go Code: Basic Concurrency

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

// Worker pool pattern
func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
    defer wg.Done()
    for job := range jobs {
        fmt.Printf("Worker %d processing job %d\n", id, job)
        time.Sleep(time.Millisecond * 100) // Simulate work
        results <- job * 2
    }
}

func main() {
    const numWorkers = 3
    const numJobs = 10

    jobs := make(chan int, numJobs)
    results := make(chan int, numJobs)
    var wg sync.WaitGroup

    // Start workers
    for w := 1; w <= numWorkers; w++ {
        wg.Add(1)
        go worker(w, jobs, results, &wg)
    }

    // Send jobs
    go func() {
        for j := 1; j <= numJobs; j++ {
            jobs <- j
        }
        close(jobs)
    }()

    // Wait for completion in another goroutine
    go func() {
        wg.Wait()
        close(results)
    }()

    // Collect results
    for result := range results {
        fmt.Printf("Result: %d\n", result)
    }
}
```

#### Go Code: Select Statement

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    ch1 := make(chan string)
    ch2 := make(chan string)

    go func() {
        time.Sleep(100 * time.Millisecond)
        ch1 <- "from channel 1"
    }()

    go func() {
        time.Sleep(200 * time.Millisecond)
        ch2 <- "from channel 2"
    }()

    // Non-blocking select with timeout
    timeout := time.After(500 * time.Millisecond)

    for i := 0; i < 2; i++ {
        select {
        case msg1 := <-ch1:
            fmt.Println("Received:", msg1)
        case msg2 := <-ch2:
            fmt.Println("Received:", msg2)
        case <-timeout:
            fmt.Println("Timeout!")
            return
        default:
            fmt.Println("No message received yet")
            time.Sleep(50 * time.Millisecond)
        }
    }
}
```

### Rust: Ownership-Based Concurrency

Rust uses ownership and borrowing rules to prevent data races at compile time.

#### Rust Code: Basic Concurrency

```rust
use std::sync::{mpsc, Arc, Mutex};
use std::thread;
use std::time::Duration;

// Worker pool pattern
fn worker(
    id: usize,
    receiver: Arc<Mutex<mpsc::Receiver<i32>>>,
    sender: mpsc::Sender<i32>,
) {
    loop {
        let job = {
            let rx = receiver.lock().unwrap();
            rx.recv()
        };

        match job {
            Ok(job) => {
                println!("Worker {} processing job {}", id, job);
                thread::sleep(Duration::from_millis(100));
                sender.send(job * 2).unwrap();
            }
            Err(_) => {
                println!("Worker {} shutting down", id);
                break;
            }
        }
    }
}

fn main() {
    const NUM_WORKERS: usize = 3;
    const NUM_JOBS: i32 = 10;

    let (job_sender, job_receiver) = mpsc::channel::<i32>();
    let (result_sender, result_receiver) = mpsc::channel::<i32>();

    let job_receiver = Arc::new(Mutex::new(job_receiver));

    // Start workers
    let mut handles = vec![];
    for w in 0..NUM_WORKERS {
        let rx = Arc::clone(&job_receiver);
        let tx = result_sender.clone();
        let handle = thread::spawn(move || worker(w, rx, tx));
        handles.push(handle);
    }

    // Send jobs
    for j in 1..=NUM_JOBS {
        job_sender.send(j).unwrap();
    }
    drop(job_sender); // Signal no more jobs
    drop(result_sender); // Drop original sender

    // Collect results
    for result in result_receiver {
        println!("Result: {}", result);
    }

    // Wait for all workers
    for handle in handles {
        handle.join().unwrap();
    }
}
```

#### Rust Code: Async/Await

```rust
use tokio::time::{sleep, Duration};
use tokio::sync::mpsc;

#[tokio::main]
async fn main() {
    let (tx1, mut rx1) = mpsc::channel(32);
    let (tx2, mut rx2) = mpsc::channel(32);

    tokio::spawn(async move {
        sleep(Duration::from_millis(100)).await;
        tx1.send("from channel 1").await.unwrap();
    });

    tokio::spawn(async move {
        sleep(Duration::from_millis(200)).await;
        tx2.send("from channel 2").await.unwrap();
    });

    let timeout = sleep(Duration::from_millis(500));
    tokio::pin!(timeout);

    for _ in 0..2 {
        tokio::select! {
            Some(msg) = rx1.recv() => {
                println!("Received: {}", msg);
            }
            Some(msg) = rx2.recv() => {
                println!("Received: {}", msg);
            }
            _ = &mut timeout => {
                println!("Timeout!");
                return;
            }
        }
    }
}
```

---

## Memory Management

### Go: Garbage Collection

```go
package main

import (
    "fmt"
    "runtime"
    "time"
)

type LargeStruct struct {
    Data [1024 * 1024]byte // 1MB
    ID   int
}

func allocateMemory() *LargeStruct {
    // Go allocates on heap automatically
    ls := &LargeStruct{ID: 42}
    return ls
}

func main() {
    var m1, m2 runtime.MemStats

    runtime.GC()
    runtime.ReadMemStats(&m1)

    // Allocate many objects
    objects := make([]*LargeStruct, 100)
    for i := range objects {
        objects[i] = allocateMemory()
    }

    runtime.ReadMemStats(&m2)

    fmt.Printf("Heap before: %d MB\n", m1.HeapAlloc/1024/1024)
    fmt.Printf("Heap after: %d MB\n", m2.HeapAlloc/1024/1024)

    // Clear references
    objects = nil

    runtime.GC()
    time.Sleep(time.Millisecond) // Allow GC to run

    runtime.ReadMemStats(&m1)
    fmt.Printf("Heap after GC: %d MB\n", m1.HeapAlloc/1024/1024)
}
```

**Go GC Characteristics:**
- Concurrent mark-and-sweep
- Low latency (sub-millisecond pauses in Go 1.20+)
- Automatic tuning
- Memory overhead ~2x live data

### Rust: Ownership System

```rust
use std::sync::Arc;

struct LargeStruct {
    data: Box<[u8; 1024 * 1024]>, // 1MB on heap
    id: i32,
}

impl LargeStruct {
    fn new(id: i32) -> Self {
        LargeStruct {
            data: Box::new([0u8; 1024 * 1024]),
            id,
        }
    }
}

fn demonstrate_ownership() {
    // Ownership transfer
    let ls1 = LargeStruct::new(1);
    let ls2 = ls1; // ls1 moved to ls2
    // println!("{}", ls1.id); // ERROR: value moved
    println!("Owner: {}", ls2.id);

    // Borrowing
    let ls3 = LargeStruct::new(3);
    borrow_struct(&ls3);
    println!("Still own: {}", ls3.id); // OK, borrowed, not moved

    // Shared ownership with Arc
    let ls4 = Arc::new(LargeStruct::new(4));
    let ls5 = Arc::clone(&ls4);
    println!("Shared owners: {} and {}", ls4.id, ls5.id);
}

fn borrow_struct(ls: &LargeStruct) {
    println!("Borrowed: {}", ls.id);
}

fn main() {
    demonstrate_ownership();
}
```

**Rust Memory Characteristics:**
- Zero-cost deterministic cleanup
- Compile-time memory safety
- RAII (Resource Acquisition Is Initialization)
- No runtime overhead

---

## Performance Characteristics

### Benchmark Comparison

| Metric | Go | Rust | Notes |
|--------|-----|------|-------|
| Binary Size | 2-5 MB | 200KB-2 MB | Rust smaller with optimization |
| Memory Usage | 2-5x live data | 1x + small overhead | Go has GC overhead |
| Startup Time | ~100ms | ~10ms | Rust faster |
| Compilation | ~1-10s | ~10-60s | Go much faster |
| Raw Speed | ~80% of C | ~98% of C | Rust closer to C++ |
| GC Pauses | 0.5-10ms | 0ms | No GC in Rust |

### HTTP Server Performance

```go
// Go HTTP Server
package main

import (
    "fmt"
    "net/http"
    "runtime"
)

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, World!")
    })

    http.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        fmt.Fprintf(w, `{"message": "Hello", "count": 42}`)
    })

    http.ListenAndServe(":8080", nil)
}
```

```rust
// Rust HTTP Server (with Actix-web)
use actix_web::{get, web, App, HttpResponse, HttpServer, Responder};
use serde::Serialize;

#[derive(Serialize)]
struct Message {
    message: String,
    count: i32,
}

#[get("/")]
async fn hello() -> impl Responder {
    "Hello, World!"
}

#[get("/json")]
async fn json_response() -> impl Responder {
    HttpResponse::Ok().json(Message {
        message: "Hello".to_string(),
        count: 42,
    })
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    HttpServer::new(|| {
        App::new()
            .service(hello)
            .service(json_response)
    })
    .bind("127.0.0.1:8080")?
    .workers(num_cpus::get())
    .run()
    .await
}
```

**Benchmark Results (approximate on modern hardware):**

| Test | Go (net/http) | Rust (Actix) | Ratio |
|------|---------------|--------------|-------|
| Hello World RPS | 150,000 | 450,000 | 3:1 |
| JSON Response RPS | 120,000 | 380,000 | 3.2:1 |
| Latency p99 | 2ms | 0.5ms | 4:1 |
| Memory @ 100k RPS | 200MB | 50MB | 4:1 |

---

## Decision Matrix

### Use Go When...

| Criterion | Weight | Go Score | Notes |
|-----------|--------|----------|-------|
| Rapid Development | High | 9/10 | Fast compile, simple syntax |
| Team Size > 10 | High | 9/10 | Easy to read/maintain |
| Microservices | High | 9/10 | Fast startup, small binaries |
| Cloud Native | High | 10/10 | Kubernetes, Docker written in Go |
| Web Services | Medium | 8/10 | Excellent standard library |
| Learning Curve | Medium | 9/10 | Learn in days, master in weeks |
| Library Ecosystem | Medium | 8/10 | Mature, well-documented |
| Hiring Pool | Medium | 8/10 | Growing rapidly |

### Use Rust When...

| Criterion | Weight | Rust Score | Notes |
|-----------|--------|------------|-------|
| Maximum Performance | High | 10/10 | Zero-cost abstractions |
| Memory Safety Critical | High | 10/10 | Compile-time guarantees |
| Systems Programming | High | 10/10 | OS kernels, embedded |
| No GC Tolerated | High | 10/10 | Real-time systems |
| Embedded/IoT | High | 9/10 | Small binaries, no runtime |
| WebAssembly | Medium | 10/10 | Primary language for WASM |
| Game Development | Medium | 9/10 | Bevy, winit ecosystem |
| Cryptography | Medium | 9/10 | Constant-time guarantees |

---

## Migration Guide

### Go to Rust Migration

#### Phase 1: Identify Components
```
1. Performance-critical paths
2. Memory-intensive operations
3. Concurrent data structures
4. Low-level system interfaces
```

#### Phase 2: FFI Bridge (Incremental Migration)

```go
// Go side: export function
package main

/*
#include <stdint.h>
extern void process_data_rust(const char* data, size_t len);
*/
import "C"
import "unsafe"

func callRust(data []byte) {
    if len(data) == 0 {
        return
    }
    C.process_data_rust(
        (*C.char)(unsafe.Pointer(&data[0])),
        C.size_t(len(data)),
    )
}
```

```rust
// Rust side: create C-compatible library
#[no_mangle]
pub extern "C" fn process_data_rust(data: *const u8, len: usize) {
    let slice = unsafe {
        std::slice::from_raw_parts(data, len)
    };
    
    // Process data in Rust
    let result = process(slice);
    
    // Return or callback to Go
}

fn process(data: &[u8]) -> Vec<u8> {
    // Implementation
    data.to_vec()
}
```

#### Phase 3: Data Structure Mapping

| Go | Rust | Notes |
|----|------|-------|
| `map[K]V` | `HashMap<K, V>` | Use `BTreeMap` for ordered |
| `[]T` | `Vec<T>` | Growable array |
| `chan T` | `mpsc::channel<T>` | Single-owner channels |
| `interface{}` | `dyn Trait` or `enum` | Prefer enums for performance |
| `*T` | `Box<T>` or `&T` | Ownership-aware |
| `sync.Mutex` | `std::sync::Mutex` | Poisoning behavior differs |
| `sync.WaitGroup` | Scoped threads or `futures::join!` | Different patterns |

### Rust to Go Migration

#### Phase 1: Component Analysis
```
1. Identify components needing GC tolerance
2. Find areas benefiting from faster compilation
3. Locate web services/API endpoints
4. Identify rapid prototyping areas
```

#### Phase 2: FFI Bridge

```rust
// Rust side: export for cgo
#[no_mangle]
pub extern "C" fn analyze_data_go(data: *const u8, len: usize) -> *mut c_char {
    let slice = unsafe {
        std::slice::from_raw_parts(data, len)
    };
    
    let result = analyze(slice);
    
    // Convert to C string (must be freed by Go)
    CString::new(result)
        .unwrap()
        .into_raw()
}
```

#### Phase 3: Common Patterns

| Rust | Go | Notes |
|------|-----|-------|
| `Result<T, E>` | `(T, error)` | Go uses multiple returns |
| `Option<T>` | `*T` or `nil` | Go uses nil for absence |
| `match` | `switch` | Go switch is more limited |
| `impl Trait for Type` | Interface satisfaction | Go interfaces are implicit |
| `mod` | `package` | Different visibility rules |
| `cargo` | `go mod` | Similar dependency management |

---

## When to Choose Which

### Choose Go For:

1. **Cloud-Native Microservices**
   - Docker, Kubernetes ecosystems
   - Fast CI/CD cycles
   - Horizontal scaling

2. **Developer Velocity Projects**
   - Startups with time pressure
   - Prototypes and MVPs
   - Teams with mixed experience

3. **DevOps and Infrastructure**
   - CLI tools
   - Deployment automation
   - Monitoring agents

4. **Network Services**
   - API gateways
   - Proxy servers
   - Load balancers

### Choose Rust For:

1. **Performance-Critical Systems**
   - Game engines
   - Real-time trading
   - High-frequency data processing

2. **Systems Programming**
   - Operating systems
   - Device drivers
   - Embedded systems

3. **Safety-Critical Applications**
   - Medical devices
   - Aerospace
   - Cryptographic implementations

4. **Resource-Constrained Environments**
   - WebAssembly
   - IoT devices
   - Serverless functions (cold start)

---

## Summary Table

| Aspect | Go | Rust | Winner |
|--------|-----|------|--------|
| Learning Curve | Gentle | Steep | Go |
| Development Speed | Fast | Moderate | Go |
| Runtime Performance | Good | Excellent | Rust |
| Memory Safety | GC-based | Compile-time | Tie |
| Concurrency Safety | Good | Excellent | Rust |
| Binary Size | Larger | Smaller | Rust |
| Compile Time | Fast | Slow | Go |
| Ecosystem Maturity | Very High | High | Go |
| Hiring Availability | High | Medium | Go |
| Long-term Maintenance | Easy | Moderate | Go |
| Systems Programming | Good | Excellent | Rust |
| Web Services | Excellent | Good | Go |

---

## Conclusion

Both Go and Rust are excellent languages that have earned their place in modern software development. The choice between them often comes down to:

- **Go**: When you prioritize developer productivity, fast iteration, and have tolerance for GC
- **Rust**: When you need maximum performance, memory safety without GC, and can invest in learning

Many successful organizations use both: Go for microservices and web APIs, Rust for performance-critical components and systems programming.

---

*Document Version: 1.0*
*Last Updated: 2026-04-03*
*Size: ~25KB*
