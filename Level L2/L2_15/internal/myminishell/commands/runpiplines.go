package commands

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

// RunPipeline - пайп лайн команд
func RunPipeline(pipeline [][]string) error {
	var prevStdout io.ReadCloser
	var processes []*exec.Cmd

	for i, args := range pipeline {
		if len(args) == 0 {
			continue
		}

		cmd := exec.Command(args[0], args[1:]...)

		if i == 0 {
			cmd.Stdin = os.Stdin
		} else {
			cmd.Stdin = prevStdout
		}

		if i == len(pipeline)-1 {
			cmd.Stdout = os.Stdout
		} else {
			stdoutPipe, err := cmd.StdoutPipe()
			if err != nil {
				return err
			}
			prevStdout = stdoutPipe
		}

		processes = append(processes, cmd)
	}

	err := RunProcess(processes, prevStdout)
	if err != nil {
		return err
	}

	return nil
}

// RunProcess - запускаем процесс пайп лайна
func RunProcess(processes []*exec.Cmd, prevStdout io.ReadCloser) error {
	for _, p := range processes {
		if err := p.Start(); err != nil {
			return err
		}
	}

	// Закрываем все промежуточные пайпы после старта
	for i := 0; i < len(processes)-1; i++ {
		if prevStdout != nil {
			prevStdout.Close()
		}
	}

	// Ждём завершения процессов
	for _, p := range processes {
		if err := p.Wait(); err != nil {
			fmt.Fprintf(os.Stderr, "Команда %s завершилась с ошибкой: %v\n", p.Args[0], err)
		}
	}

	return nil
}
