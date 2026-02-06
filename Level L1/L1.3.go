package main

import (
	"fmt"
	"time"
)

func worker(id int, works <-chan int, results chan<- int) {
	for n := range works {
		fmt.Println("Worker", id, "start", n)
		time.Sleep(time.Millisecond * time.Duration(n))
		fmt.Println("Worker", id, "end", n)
		results <- n
	}
}

func main() {
	works := make(chan int, 50)
	results := make(chan int, 50)

	var workerCount int
	fmt.Println("Введите кол-во воркеров: ")
	if _, err := fmt.Scan(&workerCount); err != nil {
		panic("Введен некорректный символ")
	}

	// Запускаем 5 воркеров читающих результат
	for i := 1; i <= workerCount; i++ {
		go worker(i, works, results)
	}

	//выводим что прочитал воркер
	go func() {
		for n := range results {
			fmt.Println(n)
		}
	}()

	// постоянно пишем
	workID := 1
	for {
		works <- workID
		workID++
		time.Sleep(time.Millisecond * 200)
	}
}
