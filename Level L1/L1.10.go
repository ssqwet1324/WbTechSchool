package main

import (
	"fmt"
	"math"
)

func main() {
	sp := []float64{-25.4, -27.0, 13.0, 19.0, 15.5, 24.5, -21.0, 32.5}
	storage := make(map[int][]float64)

	for _, x := range sp {
		//тут делим на 10 и округляем и умножаем обратно на 10 чтобы получить ключи
		num := math.Trunc(x/10) * 10
		//тут добавляем значения по ключу в список
		storage[int(num)] = append(storage[int(num)], x)
	}

	for k, v := range storage {
		fmt.Printf("%d:%v ", k, v)
	}
}
