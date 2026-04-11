package main

import (
	"fmt"
	"os"
	"strconv"
)

func getFlag(args []string, flag string) string {
	for i, a := range args {
		if a == flag && i+1 < len(args) {
			return args[i+1]
		}
	}
	return ""
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	switch os.Args[1] {

	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Command: todo add <title> [--desc <description>] [--priority <low|medium|high>]")
			return
		}

		title := os.Args[2]
		description := getFlag(os.Args, "--desc")
		priority := getFlag(os.Args, "--priority")
		addTask(title, description, priority)
	case "list":
		filter := ""
		if len(os.Args) >= 3 {
			filter = os.Args[2]
		}
		listTasks(filter)

	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Command: todo delete <Id>")
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Invalid task id")
			return
		}

		deleteTask(id)
	case "update":
		if len(os.Args) < 3 {
			fmt.Println("Command: todo update <id> [--title <title>] [--desc <description>] [--priority <low|medium|high>]")
			return
		}

		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Invalid task id")
			return
		}

		title := getFlag(os.Args, "--title")
		description := getFlag(os.Args, "--desc")
		priority := getFlag(os.Args, "--priority")
		updateTasks(id, title, description, priority)

	case "done":
		if len(os.Args) < 3 {
			fmt.Println("Command: todo done <id>")
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Invalid task Id")
			return
		}
		toggleDone(id)

	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		printUsage()
	}
}

func printUsage() {
	fmt.Println(`Todo CLI - Task Manager

Usage:
  todo add <title> [--desc <description>] [--priority <low|medium|high>]
  todo list [done|pending]
  todo done <id>
  todo delete <id>
  todo update <id> [--title <title>] [--desc <description>] [--priority <low|medium|high>]`)
}
