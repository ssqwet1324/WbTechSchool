package or_test

import (
	"fmt"
	"time"

	"L4.1/or"
)

// ExampleOr — пример работы функции Or
func ExampleOr() {
	sig := func(after time.Duration) <-chan interface{} {
		ch := make(chan interface{})
		go func() {
			defer close(ch)
			time.Sleep(after)
		}()
		return ch
	}

	<-or.Or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Minute),
		sig(500*time.Millisecond),
	)

	fmt.Println("done")

	// Output:
	// done
}
