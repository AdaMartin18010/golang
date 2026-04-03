# Concurrency Models Comparison

## Executive Summary

Concurrency models define how languages handle multiple simultaneous operations. This document compares CSP (Communicating Sequential Processes) in Go, Actor Model in Erlang/Elixir, Async/Await in JavaScript/C#, Threads in Java, and Coroutines in Kotlin/Rust.

---

## Table of Contents

- [Concurrency Models Comparison](#concurrency-models-comparison)
  - [Executive Summary](#executive-summary)
  - [Table of Contents](#table-of-contents)
  - [Go: CSP and Goroutines](#go-csp-and-goroutines)
  - [Erlang/Elixir: Actor Model](#erlangelixir-actor-model)
  - [JavaScript/TypeScript: Event Loop](#javascripttypescript-event-loop)
  - [Java: Threads and Executors](#java-threads-and-executors)
  - [C#: Async/Await and TPL](#c-asyncawait-and-tpl)
  - [Rust: Ownership-Based](#rust-ownership-based)
  - [Kotlin: Coroutines](#kotlin-coroutines)
  - [Comparison Matrix](#comparison-matrix)

---

## Go: CSP and Goroutines

Go implements Hoare's CSP model with goroutines and channels:

```go
package main

import (
    "context"
    "fmt"
    "sync"
    "time"
)

// Goroutines are lightweight threads
func basicGoroutines() {
    // Launch goroutine
    go func() {
        fmt.Println("Hello from goroutine")
    }()

    time.Sleep(time.Millisecond) // Wait for goroutine
}

// Channels for communication
type Message struct {
    ID   int
    Text string
}

func channelExample() {
    ch := make(chan Message, 10) // Buffered channel

    // Producer
    go func() {
        for i := 0; i < 5; i++ {
            ch <- Message{ID: i, Text: fmt.Sprintf("Message %d", i)}
        }
        close(ch)
    }()

    // Consumer
    for msg := range ch {
        fmt.Printf("Received: %+v\n", msg)
    }
}

// Select for multiplexing
func selectExample() {
    ch1 := make(chan string)
    ch2 := make(chan string)

    go func() {
        time.Sleep(100 * time.Millisecond)
        ch1 <- "from ch1"
    }()

    go func() {
        time.Sleep(200 * time.Millisecond)
        ch2 <- "from ch2"
    }()

    timeout := time.After(300 * time.Millisecond)

    for i := 0; i < 2; i++ {
        select {
        case msg := <-ch1:
            fmt.Println("Received:", msg)
        case msg := <-ch2:
            fmt.Println("Received:", msg)
        case <-timeout:
            fmt.Println("Timeout!")
            return
        }
    }
}

// Worker pool pattern
func workerPool() {
    const numWorkers = 3
    jobs := make(chan int, 100)
    results := make(chan int, 100)

    var wg sync.WaitGroup

    // Start workers
    for w := 1; w <= numWorkers; w++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            for job := range jobs {
                fmt.Printf("Worker %d processing job %d\n", id, job)
                results <- job * 2
            }
        }(w)
    }

    // Send jobs
    go func() {
        for j := 1; j <= 9; j++ {
            jobs <- j
        }
        close(jobs)
    }()

    // Close results when done
    go func() {
        wg.Wait()
        close(results)
    }()

    // Collect results
    for result := range results {
        fmt.Println("Result:", result)
    }
}

// Context for cancellation
func contextExample() {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    ch := make(chan int)

    go func() {
        for i := 0; ; i++ {
            select {
            case ch <- i:
                time.Sleep(100 * time.Millisecond)
            case <-ctx.Done():
                fmt.Println("Worker cancelled:", ctx.Err())
                return
            }
        }
    }()

    for {
        select {
        case v := <-ch:
            fmt.Println("Received:", v)
        case <-ctx.Done():
            fmt.Println("Main cancelled:", ctx.Err())
            return
        }
    }
}

// Sync primitives
type SafeCounter struct {
    mu    sync.RWMutex
    value int
}

func (c *SafeCounter) Increment() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.value++
}

func (c *SafeCounter) Value() int {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return c.value
}
```

---

## Erlang/Elixir: Actor Model

```elixir
# Elixir: Actor model with processes

# Spawn processes
pid = spawn(fn ->
  receive do
    {:ping, sender} ->
      send(sender, :pong)
  end
end)

send(pid, {:ping, self()})

receive do
  :pong -> IO.puts("Got pong!")
after
  5000 -> IO.puts("Timeout")
end

# GenServer for stateful processes
defmodule Counter do
  use GenServer

  # Client API
  def start_link(initial_value \\ 0) do
    GenServer.start_link(__MODULE__, initial_value, name: __MODULE__)
  end

  def increment(amount \\ 1) do
    GenServer.cast(__MODULE__, {:increment, amount})
  end

  def get_value do
    GenServer.call(__MODULE__, :get_value)
  end

  # Server callbacks
  @impl true
  def init(initial_value) do
    {:ok, initial_value}
  end

  @impl true
  def handle_cast({:increment, amount}, state) do
    {:noreply, state + amount}
  end

  @impl true
  def handle_call(:get_value, _from, state) do
    {:reply, state, state}
  end
end

# Supervisor for fault tolerance
defmodule Counter.Supervisor do
  use Supervisor

  def start_link(init_arg) do
    Supervisor.start_link(__MODULE__, init_arg, name: __MODULE__)
  end

  @impl true
  def init(_init_arg) do
    children = [
      Counter
    ]

    Supervisor.init(children, strategy: :one_for_one)
  end
end

# Task for async operations
task = Task.async(fn ->
  # Long computation
  :timer.sleep(1000)
  42
end)

result = Task.await(task, 5000)

# Agent for simple state
{:ok, agent} = Agent.start_link(fn -> %{} end)

Agent.update(agent, fn state ->
  Map.put(state, :key, "value")
end)

value = Agent.get(agent, fn state -> state[:key] end)

# Registry for process discovery
Registry.start_link(keys: :unique, name: MyRegistry)

{:ok, pid} = GenServer.start_link(MyWorker, [], name: {:via, Registry, {MyRegistry, "worker1"}})
pid = Registry.lookup(MyRegistry, "worker1")
```

---

## JavaScript/TypeScript: Event Loop

```typescript
// JavaScript: Event loop with async/await

// Promises
function fetchData(url: string): Promise<Response> {
    return fetch(url)
        .then(response => {
            if (!response.ok) {
                throw new Error(`HTTP ${response.status}`);
            }
            return response.json();
        });
}

// Async/await
async function getUser(id: string): Promise<User> {
    const response = await fetch(`/api/users/${id}`);
    if (!response.ok) {
        throw new Error(`Failed to fetch user: ${response.status}`);
    }
    return response.json();
}

// Parallel execution
async function getDashboardData(userId: string): Promise<Dashboard> {
    const [user, orders, notifications] = await Promise.all([
        getUser(userId),
        getOrders(userId),
        getNotifications(userId)
    ]);

    return { user, orders, notifications };
}

// Promise.race for timeouts
async function fetchWithTimeout(url: string, timeout: number): Promise<Response> {
    const timeoutPromise = new Promise<never>((_, reject) => {
        setTimeout(() => reject(new Error('Timeout')), timeout);
    });

    return Promise.race([fetch(url), timeoutPromise]);
}

// EventEmitter pattern
import { EventEmitter } from 'events';

class DataProcessor extends EventEmitter {
    async processLargeDataset(dataset: any[]) {
        for (let i = 0; i < dataset.length; i++) {
            this.emit('progress', { current: i, total: dataset.length });
            await this.processItem(dataset[i]);
        }
        this.emit('complete');
    }
}

const processor = new DataProcessor();
processor.on('progress', (info) => {
    console.log(`Progress: ${info.current}/${info.total}`);
});
processor.on('complete', () => {
    console.log('Processing complete');
});

// Worker threads for CPU-intensive tasks
import { Worker, isMainThread, parentPort, workerData } from 'worker_threads';

if (isMainThread) {
    // Main thread
    const worker = new Worker(__filename, {
        workerData: { numbers: [1, 2, 3, 4, 5] }
    });

    worker.on('message', (result) => {
        console.log('Result:', result);
    });

    worker.on('error', (err) => {
        console.error('Worker error:', err);
    });
} else {
    // Worker thread
    const { numbers } = workerData;
    const sum = numbers.reduce((a, b) => a + b, 0);
    parentPort?.postMessage(sum);
}
```

---

## Java: Threads and Executors

```java
// Java: Thread-based concurrency

// Basic thread
Thread thread = new Thread(() -> {
    System.out.println("Running in thread: " + Thread.currentThread().getName());
});
thread.start();

// Executor framework
ExecutorService executor = Executors.newFixedThreadPool(4);

// Submit tasks
Future<Integer> future = executor.submit(() -> {
    Thread.sleep(1000);
    return 42;
});

// Get result (blocks)
try {
    Integer result = future.get(5, TimeUnit.SECONDS);
} catch (InterruptedException | ExecutionException | TimeoutException e) {
    e.printStackTrace();
}

// CompletableFuture for async programming
CompletableFuture<String> future1 = CompletableFuture.supplyAsync(() -> {
    return fetchData("url1");
});

CompletableFuture<String> future2 = CompletableFuture.supplyAsync(() -> {
    return fetchData("url2");
});

// Combine futures
CompletableFuture<String> combined = future1
    .thenCombine(future2, (result1, result2) -> result1 + result2);

// Chain operations
CompletableFuture<User> userFuture = CompletableFuture
    .supplyAsync(() -> fetchUserData(userId))
    .thenApply(data -> parseUser(data))
    .thenCompose(user -> enrichUserData(user))
    .exceptionally(ex -> {
        log.error("Error fetching user", ex);
        return defaultUser();
    });

// Parallel streams
List<Integer> numbers = Arrays.asList(1, 2, 3, 4, 5, 6, 7, 8, 9, 10);
int sum = numbers.parallelStream()
    .mapToInt(n -> n * n)
    .filter(n -> n > 10)
    .sum();

// Virtual threads (Project Loom)
try (var executor = Executors.newVirtualThreadPerTaskExecutor()) {
    IntStream.range(0, 10_000).forEach(i -> {
        executor.submit(() -> {
            Thread.sleep(Duration.ofSeconds(1));
            return i;
        });
    });
}

// Structured concurrency (preview)
try (var scope = new StructuredTaskScope.ShutdownOnFailure()) {
    Future<String> user = scope.fork(() -> fetchUser());
    Future<Integer> order = scope.fork(() -> fetchOrder());

    scope.join();           // Wait for both
    scope.throwIfFailed();  // Propagate errors

    return new Response(user.resultNow(), order.resultNow());
}
```

---

## C#: Async/Await and TPL

```csharp
// C#: Task Parallel Library and async/await

// Async method
public async Task<User> GetUserAsync(string id)
{
    var response = await httpClient.GetAsync($"/api/users/{id}");
    response.EnsureSuccessStatusCode();
    return await response.Content.ReadFromJsonAsync<User>();
}

// Parallel tasks
public async Task<Dashboard> GetDashboardAsync(string userId)
{
    var userTask = GetUserAsync(userId);
    var ordersTask = GetOrdersAsync(userId);
    var statsTask = GetStatsAsync(userId);

    await Task.WhenAll(userTask, ordersTask, statsTask);

    return new Dashboard
    {
        User = await userTask,
        Orders = await ordersTask,
        Stats = await statsTask
    };
}

// Task with timeout
public async Task<T> WithTimeout<T>(Task<T> task, int milliseconds)
{
    using var cts = new CancellationTokenSource(milliseconds);

    var completed = await Task.WhenAny(task, Task.Delay(milliseconds, cts.Token));

    if (completed == task)
    {
        cts.Cancel(); // Cancel timeout task
        return await task;
    }

    throw new TimeoutException();
}

// Parallel LINQ
var result = numbers
    .AsParallel()
    .Where(n => n > 10)
    .Select(n => n * n)
    .OrderBy(n => n)
    .ToList();

// Channels (similar to Go)
Channel<int> channel = Channel.CreateUnbounded<int>();

// Producer
_ = Task.Run(async () =>
{
    for (int i = 0; i < 100; i++)
    {
        await channel.Writer.WriteAsync(i);
    }
    channel.Writer.Complete();
});

// Consumer
await foreach (var item in channel.Reader.ReadAllAsync())
{
    Console.WriteLine(item);
}

// Dataflow (TPL Dataflow)
var block = new TransformBlock<int, int>(n => n * n,
    new ExecutionDataflowBlockOptions
    {
        MaxDegreeOfParallelism = 4
    });

var actionBlock = new ActionBlock<int>(n => Console.WriteLine(n));

block.LinkTo(actionBlock);

for (int i = 0; i < 100; i++)
{
    block.Post(i);
}

block.Complete();
await actionBlock.Completion;
```

---

## Rust: Ownership-Based

```rust
// Rust: Ownership and borrowing for safe concurrency

use std::thread;
use std::sync::{mpsc, Arc, Mutex};
use std::time::Duration;

// Spawn threads
let handle = thread::spawn(|| {
    println!("Hello from thread!");
    42
});

let result = handle.join().unwrap();

// Move data into thread
let data = vec![1, 2, 3];
thread::spawn(move || {
    println!("{:?}", data);
}).join().unwrap();

// Channels for communication
let (tx, rx) = mpsc::channel();

thread::spawn(move || {
    tx.send("Hello").unwrap();
});

let message = rx.recv().unwrap();

// Shared state with Arc and Mutex
let counter = Arc::new(Mutex::new(0));
let mut handles = vec![];

for _ in 0..10 {
    let counter = Arc::clone(&counter);
    let handle = thread::spawn(move || {
        let mut num = counter.lock().unwrap();
        *num += 1;
    });
    handles.push(handle);
}

for handle in handles {
    handle.join().unwrap();
}

println!("Result: {}", *counter.lock().unwrap());

// Async/await with Tokio
#[tokio::main]
async fn main() {
    let result = fetch_data("https://example.com").await;

    // Spawn async task
    let handle = tokio::spawn(async {
        do_work().await
    });

    let result = handle.await.unwrap();

    // Select!
    tokio::select! {
        result = task1() => println!("Task 1: {:?}", result),
        result = task2() => println!("Task 2: {:?}", result),
        _ = tokio::time::sleep(Duration::from_secs(5)) => {
            println!("Timeout!");
        }
    }
}

async fn fetch_data(url: &str) -> Result<String, reqwest::Error> {
    reqwest::get(url).await?.text().await
}

// Rayon for data parallelism
use rayon::prelude::*;

let sum: i32 = (0..1_000_000)
    .into_par_iter()
    .map(|x| x * x)
    .sum();
```

---

## Kotlin: Coroutines

```kotlin
// Kotlin: Coroutines for async programming

import kotlinx.coroutines.*
import kotlinx.coroutines.channels.Channel

// Basic coroutine
fun main() = runBlocking {
    launch {
        delay(1000L)
        println("World!")
    }
    println("Hello,")
}

// Async/await
suspend fun fetchUser(id: String): User {
    delay(100) // Simulate network
    return User(id, "John")
}

suspend fun fetchDashboard(userId: String): Dashboard = coroutineScope {
    val user = async { fetchUser(userId) }
    val orders = async { fetchOrders(userId) }
    val stats = async { fetchStats(userId) }

    Dashboard(user.await(), orders.await(), stats.await())
}

// Channels (CSP style)
fun CoroutineScope.producer(): ReceiveChannel<Int> = produce {
    for (x in 1..5) {
        send(x * x)
    }
}

fun CoroutineScope.launchProcessor(id: Int, channel: ReceiveChannel<Int>) = launch {
    for (msg in channel) {
        println("Processor $id received $msg")
    }
}

// Flow for reactive streams
fun fetchUsers(): Flow<User> = flow {
    for (i in 1..3) {
        delay(100)
        emit(fetchUser(i.toString()))
    }
}.map { user ->
    user.copy(name = user.name.uppercase())
}.filter { user ->
    user.name.isNotEmpty()
}

// Usage
runBlocking {
    fetchUsers().collect { user ->
        println(user)
    }
}

// Structured concurrency
suspend fun loadData(): String = supervisorScope {
    val deferred1 = async { fetchData1() }
    val deferred2 = async { fetchData2() }

    try {
        "${deferred1.await()} ${deferred2.await()}"
    } catch (e: Exception) {
        "Error: ${e.message}"
    }
}

// Mutex for shared state
val mutex = Mutex()
var counter = 0

suspend fun increment() {
    mutex.withLock {
        counter++
    }
}
```

---

## Comparison Matrix

| Feature | Go | Erlang | JS/TS | Java | C# | Rust | Kotlin |
|---------|-----|--------|-------|------|-----|------|--------|
| Model | CSP | Actor | Event Loop | Threads | TPL | Ownership | Coroutines |
| Lightweight | Yes (2KB) | Yes (300B) | Yes | No | Yes (async) | Yes | Yes |
| Compile-time Safe | Partial | Yes | No | Yes | Yes | Yes | Yes |
| Data Race Detection | Race detector | No (isolated) | No | No | No | Yes | No |
| Cancellation | Context | Exit signal | AbortController | Interrupt | CancellationToken | Drop | Job cancel |
| Backpressure | Buffer sizing | Mailbox | Stream | Blocking | Bounded channel | Channel | Buffer |
| Debugging | Good | Excellent | Good | Hard | Good | Good | Good |

---

*Document Version: 1.0*
*Last Updated: 2026-04-03*
*Size: ~25KB*
