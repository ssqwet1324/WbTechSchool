package commands

import (
	"errors"
	"os"
)

// Pwd - текущие запущенные процессы
func Pwd() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", errors.New("cannot get current directory")
	}

	return currentDir, nil
}
