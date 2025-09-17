package commands

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// ChangeDir - сменить папку
func ChangeDir(args []string) error {
	fullPath := strings.Join(args[1:], " ")

	err := os.Chdir(fullPath)
	if err != nil {
		return errors.New("Error changing directory: " + err.Error())
	}
	fmt.Printf("Changing dir: %v\n", fullPath)

	return nil
}
