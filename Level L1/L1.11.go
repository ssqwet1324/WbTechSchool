package main

import "fmt"

// Intersection - пересечение
func Intersection(a, b []int) []int {
	set := make(map[int]struct{}, len(a))
	for _, v := range a {
		set[v] = struct{}{}
	}

	var res []int
	for _, v := range b {
		if _, ok := set[v]; ok {
			res = append(res, v)
		}
	}
	return res
}

func main() {
	sp1 := []int{1, 2, 3}
	sp2 := []int{2, 3, 4}
	fmt.Println(Intersection(sp1, sp2))

	a := []int{10, 20, 30, 40}
	b := []int{5, 10, 40, 50, 10}
	fmt.Println(Intersection(a, b))

	x := []int{7, 8, 9}
	y := []int{1, 2, 3}
	fmt.Println(Intersection(x, y))
}
