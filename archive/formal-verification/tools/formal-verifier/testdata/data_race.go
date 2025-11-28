package main

import (
	"sync"
	"time"
)

// 示例1：数据竞争 - 并发读写共享变量
var counter int

func raceExample1() {
	go func() {
		counter++
	}()

	go func() {
		counter++
	}()

	time.Sleep(100 * time.Millisecond)
}

// 示例2：数据竞争 - Map并发读写
var m = make(map[int]int)

func raceExample2() {
	go func() {
		m[1] = 1
	}()

	go func() {
		_ = m[1]
	}()

	time.Sleep(100 * time.Millisecond)
}

// 示例3：正确示例 - 使用Mutex保护
var (
	safeCounter int
	mu          sync.Mutex
)

func noRaceExample1() {
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		mu.Lock()
		safeCounter++
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		mu.Lock()
		safeCounter++
		mu.Unlock()
	}()

	wg.Wait()
}

// 示例4：正确示例 - 使用Channel通信
func noRaceExample2() {
	ch := make(chan int, 1)

	go func() {
		ch <- 42
	}()

	value := <-ch
	_ = value
}

// 示例5：复杂数据竞争 - 闭包捕获
func raceExample3() {
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// 竞争：多个goroutine读取同一个变量i
			_ = i
		}()
	}

	wg.Wait()
}

// 示例6：正确示例 - 闭包参数传递
func noRaceExample3() {
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			_ = id
		}(i)
	}

	wg.Wait()
}

func main() {
	raceExample1()
	raceExample2()
	noRaceExample1()
	noRaceExample2()
	raceExample3()
	noRaceExample3()
}
