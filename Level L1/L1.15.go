package main

var justString string

func someFunc() {
	v := createHugeString(1 << 10)
	justString = string([]rune(v[:100])) // то как должно быть реализовано
}

func main() {
	someFunc()
}
