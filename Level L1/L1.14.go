package main

import (
	"fmt"
	"reflect"
)

// CheckTypeFirstVersion - первый способ через банальный свич по всем типам
func CheckTypeFirstVersion(variable interface{}) {
	switch variable.(type) {
	case int:
		fmt.Println("int")
	case string:
		fmt.Println("string")
	case bool:
		fmt.Println("bool")
	case chan int:
		fmt.Println("chan int")
	case chan string:
		fmt.Println("chan string")
	case chan bool:
		fmt.Println("chan bool")
	default:
		fmt.Println("unknown type")
	}
}

// CheckTypeSecondVersion - 2й способ определяем канал череззз reflect
func CheckTypeSecondVersion(variable interface{}) {
	switch variable.(type) {
	case int:
		fmt.Println("int")
	case string:
		fmt.Println("string")
	case bool:
		fmt.Println("bool")
	default:
		if reflect.TypeOf(variable).Kind() == reflect.Chan {
			fmt.Println("chan")
		} else {
			fmt.Println("unknown")
		}
	}
}

// makeVals - тест значений
func makeVals() []interface{} {
	num := 1
	str := "привет"
	a := true
	chan1 := make(chan int)
	chan2 := make(chan string)
	chan3 := make(chan bool)
	b := struct{}{}

	return []interface{}{num, str, a, chan1, chan2, chan3, b}
}

func main() {
	vals := makeVals()

	fmt.Println("Первый способ")
	for _, v := range vals {
		CheckTypeFirstVersion(v)
	}

	fmt.Println("Второй способ")
	for _, v := range vals {
		CheckTypeSecondVersion(v)
	}
}
