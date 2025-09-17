package commands

import (
	"errors"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// Ps - вызов внешней команды
func Ps(args []string) error {
	systemName := strings.ToLower(runtime.GOOS)

	switch systemName {
	case "windows":
		// для винды
		cmd := exec.Command("powershell", args[0])
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return errors.New("ошибка вывода списка запущенных процессов")
		}
	case "linux":
		cmd := exec.Command("ps", "aux")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return errors.New("ошибка вывода списка запущенных процессов")
		}
	default:
		cmd := exec.Command("ps", "aux")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return errors.New("ошибка вывода списка запущенных процессов")
		}
	}

	return nil
}
