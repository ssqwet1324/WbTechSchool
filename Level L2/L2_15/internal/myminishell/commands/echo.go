package commands

import (
	"fmt"
	"strings"
)

// Echo - вывод аргументов
func Echo(args []string) {
	if len(args) < 2 {
		fmt.Println()
		return
	}

	fmt.Println(strings.Join(args[1:], " "))
}
