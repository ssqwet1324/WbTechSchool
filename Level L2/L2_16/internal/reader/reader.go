package reader

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

// ReadCommand - читаем ввод с консоли
func ReadCommand() (string, string, error) {
	var command []string
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", "", errors.New("Reader: error reading input" + err.Error())
	}

	command = strings.Fields(line)
	if command[0] != "wget" {
		return "", "", errors.New("invalid command")
	}
	if len(command) < 2 {
		return "", "", errors.New("command takes 2 arguments: wget <URL> <download depth>")
	}

	return command[1], command[2], nil
}
