package main

import (
	"fmt"
	"os"
	"time"

	"github.com/beevik/ntp"
)

func GiveNowTime() (time.Time, error) {
	t, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func main() {
	t, err := GiveNowTime()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	fmt.Println("Current time:", t)
}
