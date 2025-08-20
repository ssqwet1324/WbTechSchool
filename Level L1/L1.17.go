package main

import "fmt"

func binarySearch(nums []int, search int) int {
	left, right := 0, len(nums)-1

	for left <= right {
		mid := left + (right-left)/2

		if nums[mid] == search {
			return mid
		} else if nums[mid] > search {
			right = mid - 1
		} else {
			left = mid + 1
		}
	}

	return -1
}

func main() {
	//массив подающийся должен быть отсортирован
	fmt.Println(binarySearch([]int{1, 2, 3, 4, 5, 6, 7, 8}, 10))
}
