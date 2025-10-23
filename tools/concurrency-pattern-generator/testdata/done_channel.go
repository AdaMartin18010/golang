package main

// DoneChannel done channel模式
func DoneChannel() {
	done := make(chan struct{})
	
	go func() {
		// Do work
		for {
			select {
			case <-done:
				return
			default:
				// Continue working
			}
		}
	}()
	
	// Signal done
	close(done)
}
