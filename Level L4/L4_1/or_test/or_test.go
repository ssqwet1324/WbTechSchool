package or_test

import (
	"testing"
	"time"

	"L4.1/or"
)

// TestOr - функция для теста
func TestOr(t *testing.T) {
	sig := func(after time.Duration) <-chan interface{} {
		ch := make(chan interface{})
		go func() {
			defer close(ch)
			time.Sleep(after)
		}()
		return ch
	}

	start := time.Now()

	<-or.Or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(2*time.Second),
	)

	elapsed := time.Since(start)

	if elapsed > 1500*time.Millisecond {
		t.Fatalf("OrTest: elapsed took too long: %v", elapsed)
	}
}
