package main

import "fmt"

func test() (x int) {
	defer func() {
		x++
	}()
	x = 1
	return
}

//данная функция выведет 2, т.к тут именованный параметр и т.к в return не указана переменная,
//то defer прибавит к ней один. В итоге будет 2

func anotherTest() int {
	var x int
	defer func() {
		x++
	}()
	x = 1
	return x
}

// тут параметр итоговый не именованный, и x локальная переменная,
// из-за этого defer увеличивает только локальную x, а возвращаемое значение уже скопировано,
// поэтому возвращается значение 1

func main() {
	fmt.Println(test())
	fmt.Println(anotherTest())
	//вывод 2 1
}
