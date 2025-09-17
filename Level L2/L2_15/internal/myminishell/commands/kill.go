package commands

import (
	"errors"
	"os"
	"strconv"
)

// KillProcess - отправляет сигнал завершения процессу с указанным PID
func KillProcess(pid string) error {
	pidInt, err := strconv.Atoi(pid)
	if err != nil {
		return errors.New("invalid PID")
	}

	// проверяем существует ли такой процесс
	proc, err := os.FindProcess(pidInt)
	if err != nil {
		return errors.New("process not found")
	}

	// Завершаем процесс
	if err := proc.Kill(); err != nil {
		return errors.New("failed to kill process: " + err.Error())
	}

	return nil
}
