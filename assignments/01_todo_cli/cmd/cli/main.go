package main

import (
	"bufio"
	"fmt"
	"os"

	todo "github.com/trungnc90/learn-golang/assignments/01_todo_cli"
	"github.com/trungnc90/learn-golang/assignments/01_todo_cli/infra"
)

func main() {
	fs := infra.NewFileStore("tasks.json")
	manager := todo.New(todo.WithStorer(fs))
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Todo CLI>")
		if !scanner.Scan() {
			continue
		}

		cmd, err := todo.ParseCommand(scanner.Text())
		if err != nil {
			fmt.Println(err)
			continue
		}

		switch {
		case cmd.Exit:
			fmt.Println("Bye")
			return
		case cmd.Help:
			printUsage()
		case cmd.Add != nil:
			if err := manager.AddTask(cmd.Add); err != nil {
				fmt.Println("Error:", err)
			}
		case cmd.List != nil:
			if err := manager.ListTasks(cmd.List); err != nil {
				fmt.Println("Error:", err)
			}
		case cmd.Delete != nil:
			if err := manager.DeleteTask(cmd.Delete); err != nil {
				fmt.Println("Error:", err)
			}
		case cmd.Update != nil:
			if err := manager.UpdateTasks(cmd.Update); err != nil {
				fmt.Println("Error:", err)
			}
		case cmd.Done != nil:
			if err := manager.ToggleDone(cmd.Done); err != nil {
				fmt.Println("Error:", err)
			}
		}
	}
}

func printUsage() {
	fmt.Println(`Todo CLI - Task Manager

Commands:
  add <title> [--desc <description>] [--priority <low|medium|high>]
  list [done|pending]
  done <id>
  delete <id>
  update <id> [--title <title>] [--desc <description>] [--priority <low|medium|high>]
  help
  exit`)
}
