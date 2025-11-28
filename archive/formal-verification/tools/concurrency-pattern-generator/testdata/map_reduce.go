package main

import "sync"

// MapFunc map函数类型
type MapFunc func(interface{}) interface{}

// ReduceFunc reduce函数类型
type ReduceFunc func(interface{}, interface{}) interface{}

// MapReduce MapReduce模式
func MapReduce(data []interface{}, mapper MapFunc, reducer ReduceFunc, initial interface{}) interface{} {
	// Map phase
	mapped := make([]interface{}, len(data))
	var wg sync.WaitGroup
	
	for i, item := range data {
		wg.Add(1)
		go func(idx int, val interface{}) {
			defer wg.Done()
			mapped[idx] = mapper(val)
		}(i, item)
	}
	wg.Wait()
	
	// Reduce phase
	result := initial
	for _, item := range mapped {
		result = reducer(result, item)
	}
	
	return result
}
