package main

import (
	"bufio"
	"fmt"
	"os"

	todo "github.com/trungnc90/learn-golang/assignments/01_todo_cli"
)

func main() {
	manager := todo.New(todo.WithFileStorer("tasks.json"))
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
			if err := todo.AddTask(manager.Storer, cmd.Add); err != nil {
				fmt.Println("Error:", err)
			}
		case cmd.List != nil:
			if err := todo.ListTasks(manager.Storer, cmd.List); err != nil {
				fmt.Println("Error:", err)
			}
		case cmd.Delete != nil:
			if err := todo.DeleteTask(manager.Storer, cmd.Delete); err != nil {
				fmt.Println("Error:", err)
			}
		case cmd.Update != nil:
			if err := todo.UpdateTasks(manager.Storer, cmd.Update); err != nil {
				fmt.Println("Error:", err)
			}
		case cmd.Done != nil:
			if err := todo.ToggleDone(manager.Storer, cmd.Done); err != nil {
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
