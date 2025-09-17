package main

import (
	"L2_15/internal/myminishell/commands"
	"L2_15/internal/reader"
	"fmt"
)

func main() {
	for {
		args, eof := reader.Read()

		if eof {
			fmt.Println("\nexit")
			break
		}
		if args == nil {
			continue
		}

		if len(args) > 1 {
			// если несколько команд — это пайплайн
			err := commands.RunPipeline(args)
			if err != nil {
				fmt.Println("Ошибка в пайплайне:", err)
			}
		} else {
			// иначе просто выполняем команду
			err := commands.CallCommands(args)
			if err != nil {
				fmt.Println("Команда не распознана корректно: ", err)
			}
		}
	}
}
