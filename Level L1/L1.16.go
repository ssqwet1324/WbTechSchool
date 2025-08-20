package main

import "fmt"

func quickSort(array []int) []int {
	if len(array) <= 1 {
		return array
	}

	pivot := array[0]
	var less []int
	var greater []int

	for i := 1; i < len(array); i++ {
		if array[i] < pivot {
			less = append(less, array[i])
		} else {
			greater = append(greater, array[i])
		}
	}

	less = quickSort(less)
	greater = quickSort(greater)

	return append(append(less, pivot), greater...)
}

func main() {
	nums := []int{7, 2, 1, 6, 8, 5, 3, 4}
	sortNums := quickSort(nums)
	fmt.Println(sortNums)
}
