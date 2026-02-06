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
	// массив подающийся должен быть отсортирован
	fmt.Println(binarySearch([]int{1, 2, 3, 4, 5, 6, 7, 8}, 10)) // нет
	fmt.Println(binarySearch([]int{1, 2, 3, 4, 5, 6, 7, 8}, 6))  // есть

	fmt.Println(binarySearch([]int{1, 3, 5, 7, 9}, 1)) // первый
	fmt.Println(binarySearch([]int{1, 3, 5, 7, 9}, 9)) // последний
	fmt.Println(binarySearch([]int{1, 3, 5, 7, 9}, 4)) // между

	fmt.Println(binarySearch([]int{2, 2, 2, 3, 3, 4}, 2)) // дубли
	fmt.Println(binarySearch([]int{2, 2, 2, 3, 3, 4}, 3)) // дубли

	fmt.Println(binarySearch([]int{-10, -3, 0, 5, 12}, -3)) // отрицательные
	fmt.Println(binarySearch([]int{-10, -3, 0, 5, 12}, 11)) // нет

	fmt.Println(binarySearch([]int{5}, 5)) // один элемент
	fmt.Println(binarySearch([]int{}, 1))  // пустой
}
