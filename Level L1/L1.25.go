package main

import (
	"fmt"
	"time"
)

func MySleep1(duration time.Duration) {
	timer := time.NewTimer(duration)
	<-timer.C
}

func MySleep2(duration time.Duration) {
	done := make(chan bool)
	timer := time.NewTimer(duration)
	go func() {
		<-timer.C
		done <- true
	}()
	<-done
}

func MySleep3(duration time.Duration) {
	start := time.Now()
	end := start.Add(duration)

	for time.Now().Before(end) {
	}
}

func main() {
	fmt.Println("1")
	MySleep1(2 * time.Second)
	fmt.Println("2")
	MySleep2(2 * time.Second)
	fmt.Println("3")
	MySleep3(2 * time.Second)
	fmt.Println("4")
}
