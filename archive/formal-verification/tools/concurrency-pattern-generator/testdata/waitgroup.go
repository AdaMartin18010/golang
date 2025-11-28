package main

import "sync"

func ParallelTasks(tasks []func()) {
	var wg sync.WaitGroup
	for _, task := range tasks {
		wg.Add(1)
		go func(t func()) {
			defer wg.Done()
			t()
		}(task)
	}
	wg.Wait()
}
