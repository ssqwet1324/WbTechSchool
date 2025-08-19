package main

import (
	"fmt"
	"reflect"
)

func ChekType(variable interface{}) {
	switch variable.(type) {
	case int:
		fmt.Println("int")
	case string:
		fmt.Println("string")
	case bool:
		fmt.Println("bool")
		//1 способ через кейс
	//case chan int:
	//	fmt.Println("chan int")
	//case chan string:
	//	fmt.Println("chan string")
	//case chan bool:
	//	fmt.Println("chan bool")

	//2 способ через reflect
	default:
		if reflect.TypeOf(variable).Kind() == reflect.Chan {
			fmt.Println("chan")
		} else {
			fmt.Println("unknown type")
		}
	}
}

func main() {
	num := 1
	str := "привет"
	a := true
	chan1 := make(chan int)
	chan2 := make(chan string)
	chan3 := make(chan bool)
	b := struct{}{}

	ChekType(num)
	ChekType(str)
	ChekType(a)
	ChekType(chan1)
	ChekType(chan2)
	ChekType(chan3)
	ChekType(b)
}
