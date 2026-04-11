package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func getFlag(args []string, flag string) string {
	for i, a := range args {
		if a == flag && i+1 < len(args) {
			return args[i+1]
		}
	}
	return ""
}

// parseArgs splits an input line into tokens, keeping quoted strings together.
func parseArgs(line string) []string {
	var args []string
	var current strings.Builder
	inQuote := false

	for _, r := range line {
		switch {
		case r == '"':
			inQuote = !inQuote
		case r == ' ' && !inQuote:
			if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
		default:
			current.WriteRune(r)
		}
	}
	if current.Len() > 0 {
		args = append(args, current.String())
	}
	return args
}

func main() {

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Todo CLI>")
		if !scanner.Scan() {
			continue
		}

		args := parseArgs(scanner.Text())
		if len(args) == 0 {
			continue
		}

		command := args[0]
		if command == "exit" || command == "quit" {
			fmt.Print("Bye\n")
			return
		}

		switch command {
		case "add":
			if len(args) < 2 {
				fmt.Println("Command: todo add <title> [--desc <description>] [--priority <low|medium|high>]")
				continue
			}

			title := args[1]
			description := getFlag(args, "--desc")
			priority := getFlag(args, "--priority")
			addTask(title, description, priority)
		case "list":
			filter := ""
			if len(args) >= 2 {
				filter = args[1]
			}
			listTasks(filter)

		case "delete":
			if len(args) < 2 {
				fmt.Println("Command: todo delete <Id>")
				continue
			}
			id, err := strconv.Atoi(args[1])
			if err != nil {
				fmt.Println("Invalid task id")
				continue
			}

			deleteTask(id)
		case "update":
			if len(args) < 2 {
				fmt.Println("Command: todo update <id> [--title <title>] [--desc <description>] [--priority <low|medium|high>]")
				continue
			}

			id, err := strconv.Atoi(args[1])
			if err != nil {
				fmt.Println("Invalid task id")
				continue
			}

			title := getFlag(args, "--title")
			description := getFlag(args, "--desc")
			priority := getFlag(args, "--priority")
			updateTasks(id, title, description, priority)

		case "done":
			if len(args) < 2 {
				fmt.Println("Command: todo done <id>")
				continue
			}
			id, err := strconv.Atoi(args[1])
			if err != nil {
				fmt.Println("Invalid task Id")
				continue
			}
			toggleDone(id)

		case "help":
			printUsage()
		default:
			fmt.Printf("Unknown command: %s\n", command)
			printUsage()
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
