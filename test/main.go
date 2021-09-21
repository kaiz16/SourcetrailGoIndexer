package main

var g = "Hello"
var b = "Hi"

func main() {
	go func() {}()
	f()
	fib(42)
	FuncaaaaFromAnotherFile()
}

func f() {
}

func fib(n int) int {
	if n <= 1 {
		return n
	}
	return fib(n-1) + fib(n-1)
}
