package commands

import (
	"fmt"
	"os"
	"os/exec"
)

// runExternal - запустить внешнюю команду
func runExternal(args []string) error {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Ошибка выполнения команды %s: %v\n", args[0], err)
		return err
	}

	return nil
}
