package commands

import (
	"fmt"
)

// CallCommands - обрабатываем команды
func CallCommands(pipeline [][]string) error {
	if len(pipeline) == 0 {
		return nil
	}

	// если пайп лайн из нескольких команд — запускаем RunPipeline
	if len(pipeline) > 1 {
		return RunPipeline(pipeline)
	}

	// иначе просто выполняем одну команду
	args := pipeline[0]
	if len(args) == 0 {
		return nil
	}

	switch args[0] {
	case "cd":
		return runCd(args)
	case "pwd":
		return runPwd()
	case "echo":
		return runEcho(args)
	case "kill":
		return runKill(args)
	case "ps":
		return runPs(args)
	default:
		return runExternal(args) // поддержка внешних команд
	}
}

// runCd - запустить cd
func runCd(args []string) error {
	if len(args) < 2 {
		return nil
	}

	return ChangeDir(args)
}

// runPwd - запустить pwd
func runPwd() error {
	dir, err := Pwd()
	if err != nil {
		return err
	}
	println(dir)

	return nil
}

// runEcho - запустить echo
func runEcho(args []string) error {
	Echo(args)

	return nil
}

// runKill - убить процесс
func runKill(args []string) error {
	if len(args) < 2 {
		return nil
	}
	if err := KillProcess(args[1]); err != nil {
		return err
	}
	fmt.Printf("Процесс %v успешно завершен\n", args[1])

	return nil
}

// runPs - запустить ps
func runPs(args []string) error {
	return Ps(args)
}
