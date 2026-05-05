package main

type Config struct {
	Host string
	Port int
}

func makeInt() int { return 10 }

func main() {
	// Test 1: basic literal expression
	p := new(42)
	println(*p)

	// Test 2: composite literal - is this valid?
	cfg := new(Config{Host: "localhost", Port: 8080})
	println(cfg.Host)

	// Test 3: function call expression
	v := new(makeInt())
	println(*v)
}
