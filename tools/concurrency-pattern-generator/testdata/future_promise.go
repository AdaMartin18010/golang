package main

// Future 异步结果
type Future struct {
	result chan interface{}
	err    chan error
}

// NewFuture 创建Future
func NewFuture(fn func() (interface{}, error)) *Future {
	f := &Future{
		result: make(chan interface{}, 1),
		err:    make(chan error, 1),
	}
	
	go func() {
		res, err := fn()
		if err != nil {
			f.err <- err
		} else {
			f.result <- res
		}
	}()
	
	return f
}

// Get 获取结果
func (f *Future) Get() (interface{}, error) {
	select {
	case res := <-f.result:
		return res, nil
	case err := <-f.err:
		return nil, err
	}
}
