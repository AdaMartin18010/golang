package main

// ErrorChannel error channel模式
func ErrorChannel() error {
	errChan := make(chan error, 1)
	
	go func() {
		// Do work
		errChan <- nil // or error
	}()
	
	return <-errChan
}

// MultiError 收集多个错误
func MultiError(n int) []error {
	errChan := make(chan error, n)
	
	for i := 0; i < n; i++ {
		go func() {
			errChan <- nil
		}()
	}
	
	errors := make([]error, 0, n)
	for i := 0; i < n; i++ {
		if err := <-errChan; err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}
