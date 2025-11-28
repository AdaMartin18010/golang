package main

import "sync"

var (
	instance *Singleton
	once sync.Once
)

type Singleton struct{}

func GetInstance() *Singleton {
	once.Do(func() {
		instance = &Singleton{}
	})
	return instance
}
