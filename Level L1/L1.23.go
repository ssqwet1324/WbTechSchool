package main

import "fmt"

func main() {
	//1 способ через copy
	elem := 2
	sp := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	copy(sp[elem:], sp[elem+1:])
	sp = sp[:len(sp)-1]
	fmt.Println(sp)

	//2 способ через append
	sp2 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	newSp := append(sp[:elem], sp2[elem+1:]...)
	fmt.Println(newSp)
}
