# Memory Management Comparison Across Languages

## Executive Summary

Memory management approaches vary significantly across programming languages, from manual management in C/C++ to garbage collection in Java/Go to ownership systems in Rust. This document compares memory allocation strategies, garbage collection techniques, and performance characteristics.

---

## Table of Contents

- [Memory Management Comparison Across Languages](#memory-management-comparison-across-languages)
  - [Executive Summary](#executive-summary)
  - [Table of Contents](#table-of-contents)
  - [Go: Garbage Collection](#go-garbage-collection)
  - [Rust: Ownership System](#rust-ownership-system)
  - [Java: JVM Garbage Collection](#java-jvm-garbage-collection)
  - [C++: Manual Management](#c-manual-management)
  - [C: Raw Memory](#c-raw-memory)
  - [Swift: Automatic Reference Counting](#swift-automatic-reference-counting)
  - [C#: .NET Garbage Collection](#c-net-garbage-collection)
  - [Performance Comparison](#performance-comparison)

---

## Go: Garbage Collection

Go uses a non-generational concurrent tri-color mark-and-sweep collector:

```go
package main

import (
    "fmt"
    "runtime"
    "time"
)

// Memory allocation patterns
func allocationPatterns() {
    // Stack allocation (escape analysis)
    x := 42  // Likely on stack
    _ = x

    // Heap allocation (escapes to heap)
    ptr := &Point{X: 1, Y: 2}  // Escapes to heap
    _ = ptr

    // Slice allocation
    slice := make([]int, 1000)  // Heap
    _ = slice

    // Small slice may stay on stack
    small := make([]int, 10)
    _ = small
}

type Point struct {
    X, Y int
}

// Escape analysis demonstration
func stackAllocation() int {
    p := Point{X: 1, Y: 2}  // On stack - doesn't escape
    return p.X + p.Y
}

func heapAllocation() *Point {
    p := &Point{X: 1, Y: 2}  // On heap - escapes through return
    return p
}

// GC tuning
func gcTuning() {
    // Force GC
    runtime.GC()

    // Read GC statistics
    var m runtime.MemStats
    runtime.ReadMemStats(&m)

    fmt.Printf("HeapAlloc: %d MB\n", m.HeapAlloc/1024/1024)
    fmt.Printf("TotalAlloc: %d MB\n", m.TotalAlloc/1024/1024)
    fmt.Printf("Sys: %d MB\n", m.Sys/1024/1024)
    fmt.Printf("NumGC: %d\n", m.NumGC)
    fmt.Printf("PauseNs (last): %d µs\n", m.PauseNs[(m.NumGC+255)%256]/1000)

    // Set GC target percentage (default 100)
    // GOGC=100 means GC runs when heap doubles
    debug.SetGCPercent(100)

    // Set memory limit (Go 1.19+)
    debug.SetMemoryLimit(10 * 1024 * 1024 * 1024) // 10GB
}

// Object pooling for high-frequency allocations
type Buffer struct {
    data []byte
}

var bufferPool = sync.Pool{
    New: func() interface{} {
        return &Buffer{data: make([]byte, 4096)}
    },
}

func usePooledBuffer() {
    buf := bufferPool.Get().(*Buffer)
    defer bufferPool.Put(buf)

    // Use buffer...
    copy(buf.data, []byte("data"))
}

// Manual memory with mmap for large allocations
import "syscall"

func largeAllocation(size int) ([]byte, error) {
    return syscall.Mmap(
        -1,
        0,
        size,
        syscall.PROT_READ|syscall.PROT_WRITE,
        syscall.MAP_ANON|syscall.MAP_PRIVATE,
    )
}

// Arena pattern (Go 1.20+)
func arenaPattern() {
    // Note: arena package is experimental
    // a := arena.NewArena()
    // defer a.Free()
    // ptr := arena.New[int](a)
}

// Finalizers (use sparingly!)
func finalizerExample() {
    obj := &LargeObject{}
    runtime.SetFinalizer(obj, func(o *LargeObject) {
        fmt.Println("LargeObject being collected")
    })
}

type LargeObject struct {
    data [1024 * 1024]byte
}
```

---

## Rust: Ownership System

Rust uses ownership and borrowing for compile-time memory safety:

```rust
// Ownership rules:
// 1. Each value has an owner
// 2. Only one owner at a time
// 3. Value dropped when owner goes out of scope

fn ownership_demo() {
    // String owns heap memory
    let s1 = String::from("hello");
    let s2 = s1;  // Ownership moved to s2
    // println!("{}", s1);  // ERROR: value moved
    println!("{}", s2);     // OK

    // Clone for deep copy
    let s3 = s2.clone();
    println!("{}", s2);  // OK
    println!("{}", s3);  // OK

    // Copy types (stack only)
    let x = 5;
    let y = x;  // Copy, not move
    println!("{} {}", x, y);  // Both OK
}

// Borrowing (references)
fn borrowing() {
    let s = String::from("hello");

    // Immutable borrow
    let len = calculate_length(&s);
    println!("'{}' length: {}", s, len);  // s still valid

    // Mutable borrow
    let mut s = String::from("hello");
    change(&mut s);
    println!("{}", s);  // "hello, world"
}

fn calculate_length(s: &String) -> usize {
    s.len()
}  // s goes out of scope but doesn't drop because it doesn't own

fn change(s: &mut String) {
    s.push_str(", world");
}

// Borrowing rules:
// - Any number of immutable references OR
// - One mutable reference (but not both)

fn borrow_rules() {
    let mut s = String::from("hello");

    let r1 = &s;
    let r2 = &s;
    // let r3 = &mut s;  // ERROR: cannot borrow as mutable
    println!("{} {}", r1, r2);  // r1 and r2 used here

    let r3 = &mut s;  // OK: r1 and r2 no longer used
    println!("{}", r3);
}

// Lifetimes
fn longest<'a>(x: &'a str, y: &'a str) -> &'a str {
    if x.len() > y.len() { x } else { y }
}

// Smart pointers
use std::rc::Rc;
use std::sync::Arc;
use std::cell::RefCell;

fn smart_pointers() {
    // Rc: Reference counting (single-threaded)
    let data = Rc::new(String::from("shared"));
    let data2 = Rc::clone(&data);
    let data3 = Rc::clone(&data);
    println!("Ref count: {}", Rc::strong_count(&data));  // 3

    // Arc: Atomic reference counting (thread-safe)
    let arc_data = Arc::new(String::from("thread-safe"));
    let arc_data2 = Arc::clone(&arc_data);

    std::thread::spawn(move || {
        println!("{}", arc_data2);
    });

    // RefCell: Interior mutability
    let cell = RefCell::new(5);
    *cell.borrow_mut() += 1;
    println!("{}", cell.borrow());
}

// Custom smart pointer
struct CustomBox<T>(T);

impl<T> CustomBox<T> {
    fn new(x: T) -> CustomBox<T> {
        CustomBox(x)
    }
}

impl<T> Drop for CustomBox<T> {
    fn drop(&mut self) {
        println!("Dropping CustomBox");
    }
}

// Memory layout control
#[repr(C)]  // C-compatible layout
struct CStruct {
    a: u8,
    b: u32,
    c: u16,
}

#[repr(packed)]  // No padding
struct PackedStruct {
    a: u8,
    b: u32,
    c: u16,
}

// Allocators (nightly)
use std::alloc::{alloc, dealloc, Layout};

unsafe fn manual_alloc() {
    let layout = Layout::new::<u32>();
    let ptr = alloc(layout) as *mut u32;
    *ptr = 42;
    dealloc(ptr as *mut u8, layout);
}
```

---

## Java: JVM Garbage Collection

Java offers multiple GC algorithms for different workloads:

```java
// JVM GC options:
// -XX:+UseG1GC (default since Java 9)
// -XX:+UseZGC (low latency, Java 15+)
// -XX:+UseShenandoahGC (low latency)
// -XX:MaxGCPauseMillis=200

import java.lang.management.*;
import java.util.*;

public class GCDemo {

    // Memory pools monitoring
    public static void printMemoryInfo() {
        for (MemoryPoolMXBean pool : ManagementFactory.getMemoryPoolMXBeans()) {
            MemoryUsage usage = pool.getUsage();
            System.out.printf("%s: %d MB / %d MB%n",
                pool.getName(),
                usage.getUsed() / 1024 / 1024,
                usage.getMax() / 1024 / 1024);
        }
    }

    // GC statistics
    public static void printGCInfo() {
        for (GarbageCollectorMXBean gc : ManagementFactory.getGarbageCollectorMXBeans()) {
            System.out.printf("%s: %d collections, %d ms total%n",
                gc.getName(),
                gc.getCollectionCount(),
                gc.getCollectionTime());
        }
    }

    // Soft references for cache
    public static class SoftCache<K, V> {
        private final Map<K, SoftReference<V>> cache = new HashMap<>();

        public V get(K key) {
            SoftReference<V> ref = cache.get(key);
            return ref != null ? ref.get() : null;
        }

        public void put(K key, V value) {
            cache.put(key, new SoftReference<>(value));
        }
    }

    // Weak references for listeners
    public static class WeakListenerManager {
        private final List<WeakReference<Listener>> listeners = new ArrayList<>();

        public void addListener(Listener listener) {
            listeners.add(new WeakReference<>(listener));
        }

        public void notifyListeners() {
            listeners.removeIf(ref -> ref.get() == null);
            for (WeakReference<Listener> ref : listeners) {
                Listener listener = ref.get();
                if (listener != null) {
                    listener.onEvent();
                }
            }
        }
    }

    interface Listener {
        void onEvent();
    }

    // Object pooling
    public static class ObjectPool<T> {
        private final Queue<T> pool;
        private final Factory<T> factory;

        public ObjectPool(int size, Factory<T> factory) {
            this.factory = factory;
            this.pool = new ArrayDeque<>(size);
            for (int i = 0; i < size; i++) {
                pool.offer(factory.create());
            }
        }

        public T borrow() {
            T obj = pool.poll();
            return obj != null ? obj : factory.create();
        }

        public void release(T obj) {
            pool.offer(obj);
        }

        interface Factory<T> {
            T create();
        }
    }

    // Off-heap memory with ByteBuffer
    public static void offHeapMemory() {
        ByteBuffer direct = ByteBuffer.allocateDirect(1024 * 1024);  // 1MB off-heap
        direct.putInt(0, 42);
        int value = direct.getInt(0);

        // Cleaner for explicit deallocation
        ((DirectBuffer) direct).cleaner().clean();
    }

    // Escape analysis optimization
    public static int stackAllocation() {
        // Point may be allocated on stack (escape analysis)
        Point p = new Point(1, 2);
        return p.x + p.y;
    }

    static class Point {
        int x, y;
        Point(int x, int y) { this.x = x; this.y = y; }
    }
}
```

---

## C++: Manual Management

C++ provides manual control with smart pointers for safety:

```cpp
#include <memory>
#include <vector>

// Raw pointers (manual management)
void raw_pointers() {
    int* ptr = new int(42);
    delete ptr;  // Must manually delete
}

// Unique ownership
void unique_ptr_demo() {
    std::unique_ptr<int> ptr = std::make_unique<int>(42);
    // ptr automatically deleted when out of scope

    // Transfer ownership
    std::unique_ptr<int> ptr2 = std::move(ptr);
    // ptr is now nullptr
}

// Shared ownership
void shared_ptr_demo() {
    std::shared_ptr<int> ptr1 = std::make_shared<int>(42);
    {
        std::shared_ptr<int> ptr2 = ptr1;  // Reference count = 2
        // Reference count = 1 after ptr2 scope ends
    }
    // Reference count = 0, memory freed
}

// Weak pointers to break cycles
void weak_ptr_demo() {
    std::shared_ptr<Node> node1 = std::make_shared<Node>();
    std::shared_ptr<Node> node2 = std::make_shared<Node>();

    node1->next = node2;
    node2->prev = node1;  // weak_ptr prevents cycle
}

struct Node {
    std::shared_ptr<Node> next;
    std::weak_ptr<Node> prev;  // Weak to break cycle
};

// Custom deleter
void custom_deleter() {
    FILE* file = fopen("test.txt", "r");
    std::unique_ptr<FILE, decltype(&fclose)> file_ptr(file, fclose);
    // fclose called automatically
}

// Custom allocator
template<typename T>
class PoolAllocator {
    struct Pool {
        alignas(alignof(T)) char memory[sizeof(T) * 1024];
        bool used[1024] = {};
    };

    std::vector<std::unique_ptr<Pool>> pools;

public:
    T* allocate() {
        for (auto& pool : pools) {
            for (int i = 0; i < 1024; ++i) {
                if (!pool->used[i]) {
                    pool->used[i] = true;
                    return reinterpret_cast<T*>(&pool->memory[i * sizeof(T)]);
                }
            }
        }

        // Create new pool
        pools.push_back(std::make_unique<Pool>());
        pools.back()->used[0] = true;
        return reinterpret_cast<T*>(&pools.back()->memory[0]);
    }

    void deallocate(T* p) {
        for (auto& pool : pools) {
            if (p >= reinterpret_cast<T*>(pool->memory) &&
                p < reinterpret_cast<T*>(pool->memory + sizeof(pool->memory))) {
                size_t idx = (reinterpret_cast<char*>(p) - pool->memory) / sizeof(T);
                pool->used[idx] = false;
                p->~T();  // Explicit destructor call
                return;
            }
        }
    }
};

// Memory-mapped files
#include <sys/mman.h>
#include <fcntl.h>

void memory_mapped_file() {
    int fd = open("large_file.dat", O_RDONLY);

    void* mapped = mmap(nullptr, file_size, PROT_READ, MAP_PRIVATE, fd, 0);

    // Access memory directly
    int value = *reinterpret_cast<int*>(mapped);

    munmap(mapped, file_size);
    close(fd);
}

// Placement new
void placement_new() {
    alignas(int) char buffer[sizeof(int)];
    int* p = new (buffer) int(42);  // Construct in existing memory
    p->~int();  // Explicit destructor
}
```

---

## C: Raw Memory

C provides direct memory control:

```c
#include <stdlib.h>
#include <string.h>

// Stack allocation
void stack_allocation() {
    int arr[100];  // 400 bytes on stack
    char buffer[1024];
}

// Heap allocation
void heap_allocation() {
    int* arr = malloc(100 * sizeof(int));
    if (arr == NULL) {
        // Handle allocation failure
        return;
    }

    memset(arr, 0, 100 * sizeof(int));

    free(arr);  // Must free
}

// Realloc for resizing
void realloc_example() {
    int* arr = malloc(10 * sizeof(int));

    // Grow array
    int* new_arr = realloc(arr, 20 * sizeof(int));
    if (new_arr != NULL) {
        arr = new_arr;
    }

    free(arr);
}

// Memory pools
typedef struct Pool {
    char* memory;
    size_t size;
    size_t used;
} Pool;

Pool* pool_create(size_t size) {
    Pool* pool = malloc(sizeof(Pool));
    pool->memory = malloc(size);
    pool->size = size;
    pool->used = 0;
    return pool;
}

void* pool_alloc(Pool* pool, size_t size) {
    if (pool->used + size > pool->size) {
        return NULL;
    }
    void* ptr = pool->memory + pool->used;
    pool->used += size;
    return ptr;
}

void pool_reset(Pool* pool) {
    pool->used = 0;
}

void pool_destroy(Pool* pool) {
    free(pool->memory);
    free(pool);
}

// Custom allocators
struct Allocator {
    void* (*alloc)(struct Allocator* self, size_t size);
    void (*free)(struct Allocator* self, void* ptr);
};

// Arena allocator
struct Arena {
    struct Allocator base;
    char* memory;
    size_t size;
    size_t used;
};

void* arena_alloc(struct Allocator* self, size_t size) {
    struct Arena* arena = (struct Arena*)self;
    if (arena->used + size > arena->size) {
        return NULL;
    }
    void* ptr = arena->memory + arena->used;
    arena->used += size;
    return ptr;
}

void arena_free_all(struct Allocator* self) {
    struct Arena* arena = (struct Arena*)self;
    arena->used = 0;
}
```

---

## Swift: Automatic Reference Counting

Swift uses ARC with compile-time optimization:

```swift
import Foundation

// ARC automatically manages memory
class Person {
    let name: String
    init(name: String) {
        self.name = name
        print("\(name) is initialized")
    }
    deinit {
        print("\(name) is deinitialized")
    }
}

// Strong references (default)
var person1: Person? = Person(name: "John")  // Reference count = 1
var person2 = person1                         // Reference count = 2
person1 = nil                                 // Reference count = 1
person2 = nil                                 // Reference count = 0, deinit called

// Weak references to avoid retain cycles
class Apartment {
    let unit: String
    weak var tenant: Person?  // Weak reference

    init(unit: String) {
        self.unit = unit
        print("Apartment \(unit) is initialized")
    }
    deinit {
        print("Apartment \(unit) is deinitialized")
    }
}

// Unowned references
class Customer {
    let name: String
    var card: CreditCard?

    init(name: String) {
        self.name = name
    }
}

class CreditCard {
    let number: String
    unowned let customer: Customer  // Always has a customer

    init(number: String, customer: Customer) {
        self.number = number
        self.customer = customer
    }
}

// Value types (copied)
struct Point {
    var x: Double
    var y: Double
}

var point1 = Point(x: 0, y: 0)
var point2 = point1  // Copy, not reference
point2.x = 10
print(point1.x)  // Still 0

// Copy-on-write for collections
var array1 = [1, 2, 3]
var array2 = array1  // Shares storage
array2.append(4)     // Copies on write

// Autorelease pools
func processLargeData() {
    for i in 0..<100000 {
        autoreleasepool {
            let data = loadData(index: i)
            process(data)
            // Data released at end of autoreleasepool
        }
    }
}
```

---

## C#: .NET Garbage Collection

```csharp
using System;
using System.Runtime;

public class MemoryManagementDemo
{
    // GC configuration
    public static void ConfigureGC()
    {
        // Server GC for high-throughput apps
        // <ServerGarbageCollection>true</ServerGarbageCollection>

        // Concurrent GC (default)
        GCSettings.LatencyMode = GCLatencyMode.Interactive;

        // Low latency mode for real-time
        // GCSettings.LatencyMode = GCLatencyMode.SustainedLowLatency;
    }

    // Force collection
    public static void ForceGC()
    {
        GC.Collect();
        GC.WaitForPendingFinalizers();
        GC.Collect();
    }

    // Memory information
    public static void PrintMemoryInfo()
    {
        Console.WriteLine($"Total memory: {GC.GetTotalMemory(false) / 1024 / 1024} MB");
        Console.WriteLine($"Generation: {GC.GetGeneration(this)}");
        Console.WriteLine($"Max generation: {GC.MaxGeneration}");
    }

    // Large object heap compaction
    public static void CompactLOH()
    {
        GCSettings.LargeObjectHeapCompactionMode =
            GCLargeObjectHeapCompactionMode.CompactOnce;
        GC.Collect(2, GCCollectionMode.Default, true, true);
    }

    // Weak references
    public static void WeakReferenceDemo()
    {
        var strong = new byte[1024 * 1024];
        var weak = new WeakReference<byte[]>(strong);

        strong = null;
        GC.Collect();

        if (weak.TryGetTarget(out byte[]? target))
        {
            Console.WriteLine("Object still alive");
        }
        else
        {
            Console.WriteLine("Object collected");
        }
    }

    // Object pooling
    public class ObjectPool<T> where T : class
    {
        private readonly Stack<T> _pool;
        private readonly Func<T> _factory;

        public ObjectPool(int capacity, Func<T> factory)
        {
            _factory = factory;
            _pool = new Stack<T>(capacity);

            for (int i = 0; i < capacity; i++)
            {
                _pool.Push(factory());
            }
        }

        public T Rent()
        {
            return _pool.Count > 0 ? _pool.Pop() : _factory();
        }

        public void Return(T obj)
        {
            _pool.Push(obj);
        }
    }

    // Disposable pattern
    public class Resource : IDisposable
    {
        private bool _disposed;
        private IntPtr _nativeResource;

        public void Dispose()
        {
            Dispose(true);
            GC.SuppressFinalize(this);
        }

        protected virtual void Dispose(bool disposing)
        {
            if (!_disposed)
            {
                if (disposing)
                {
                    // Dispose managed resources
                }

                // Free native resources
                FreeNativeResource(_nativeResource);
                _disposed = true;
            }
        }

        ~Resource()
        {
            Dispose(false);
        }
    }
}
```

---

## Performance Comparison

| Aspect | Go | Rust | Java | C++ | C | Swift | C# |
|--------|-----|------|------|-----|---|-------|-----|
| Allocation Speed | Fast | Fast | Moderate | Fast | Fast | Fast | Moderate |
| Deallocation | GC (pause) | Deterministic | GC (pause) | Deterministic | Manual | ARC | GC (pause) |
| Memory Overhead | 2-3x | 1x | 2-3x | 1x | 1x | 1x | 2-3x |
| Pause Latency | 0.5-2ms | 0ms | 5-50ms | 0ms | 0ms | 0ms | 5-50ms |
| Safety | Safe | Safe | Safe | Unsafe | Unsafe | Safe | Safe |
| Learning Curve | Easy | Hard | Easy | Hard | Hard | Easy | Easy |
| Real-time | No | Yes | No | Yes | Yes | Yes | No |

---

*Document Version: 1.0*
*Last Updated: 2026-04-03*
*Size: ~25KB*
