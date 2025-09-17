package reader

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/google/shlex"
)

// Read - читаем консоль
func Read() ([][]string, bool) {
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, true
	}

	line = strings.TrimSpace(line)
	line = strings.ReplaceAll(line, `\`, `\\`)

	// Разбиваем по |
	parts := strings.Split(line, "|")
	commands := make([][]string, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		args, err := shlex.Split(part)
		if err != nil {
			fmt.Println("Ошибка парсинга:", err)
			return nil, false
		}
		commands = append(commands, args)
	}

	return commands, false
}
