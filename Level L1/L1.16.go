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
	nums1 := []int{7, 2, 1, 6, 8, 5, 3, 4}
	fmt.Println(quickSort(nums1))

	nums2 := []int{5, 4, 3, 2, 1}
	fmt.Println(quickSort(nums2))

	nums3 := []int{1, 2, 3, 4, 5}
	fmt.Println(quickSort(nums3))

	nums4 := []int{3, 1, 2, 3, 2, 1}
	fmt.Println(quickSort(nums4))

	nums5 := []int{42}
	fmt.Println(quickSort(nums5))

	var nums6 []int
	fmt.Println(quickSort(nums6))

	nums7 := []int{0, -10, 5, -3, 2}
	fmt.Println(quickSort(nums7))
}
