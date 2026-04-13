package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	store := NewFileStore("tasks.json")
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Todo CLI>")
		if !scanner.Scan() {
			continue
		}

		cmd, err := parseCommand(scanner.Text())
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
			addTask(store, cmd.Add)
		case cmd.List != nil:
			listTasks(store, cmd.List)
		case cmd.Delete != nil:
			deleteTask(store, cmd.Delete)
		case cmd.Update != nil:
			updateTasks(store, cmd.Update)
		case cmd.Done != nil:
			toggleDone(store, cmd.Done)
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
