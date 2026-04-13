package main

import (
	"bufio"
	"fmt"
	"os"

	todo "github.com/trungnc90/learn-golang/assignments/01_todo_cli"
)

func main() {
	store := todo.NewFileStore("tasks.json")
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
			todo.AddTask(store, cmd.Add)
		case cmd.List != nil:
			todo.ListTasks(store, cmd.List)
		case cmd.Delete != nil:
			todo.DeleteTask(store, cmd.Delete)
		case cmd.Update != nil:
			todo.UpdateTasks(store, cmd.Update)
		case cmd.Done != nil:
			todo.ToggleDone(store, cmd.Done)
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
