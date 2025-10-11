package main

import (
	"fmt"
	"time"
)

//
//import (
//	"fmt"
//	"time"
//)
//
//func writer() <-chan int {
//	c := make(chan int)
//	go func() {
//		for i := 0; i < 10; i++ {
//			c <- i
//		}
//		close(c)
//	}()
//
//	return c
//}
//
//func doubler(ch <-chan int) <-chan int {
//	c := make(chan int)
//
//	go func() {
//		for i := range ch {
//			time.Sleep(time.Millisecond * 500)
//			c <- i * 2
//		}
//		close(c)
//	}()
//
//	return c
//}
//
//func reader(ch <-chan int) {
//	for i := range ch {
//		fmt.Println(i)
//	}
//}

//func randomTimeWork() {
//	time.Sleep(time.Duration(rand.Intn(100)) * time.Second)
//}
//
//func predictableTimeWork() {
//	ch := make(chan struct{})
//
//	go func() {
//		randomTimeWork()
//		close(ch)
//	}()
//
//	select {
//	case <-ch:
//	case <-time.After(3 * time.Second):
//		panic("time out")
//	}
//}

func main() {
	now := time.Now()
	now.Date()
	fmt.Println(now)

}
