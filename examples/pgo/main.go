package main

import (
	"fmt"
	"math/rand"
	"time"
)

// 热点函数，用于演示 PGO 优化内联热点
func hotPath(data []int) int {
	sum := 0
	for _, v := range data {
		// 简单热点运算
		sum += (v * 3) ^ 7
	}
	return sum
}

func main() {
	rand.Seed(42)
	n := 1_0000
	data := make([]int, n)
	for i := 0; i < n; i++ {
		data[i] = rand.Intn(1000)
	}
	start := time.Now()
	total := 0
	for i := 0; i < 2000; i++ {
		total ^= hotPath(data)
	}
	fmt.Println("result:", total, "elapsed:", time.Since(start))
}
