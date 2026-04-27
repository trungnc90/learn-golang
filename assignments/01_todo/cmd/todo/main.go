package main

import (
	"bufio"
	"fmt"
	"os"
	"text/tabwriter"

	todo "github.com/trungnc90/learn-golang/assignments/01_todo"
	"github.com/trungnc90/learn-golang/assignments/01_todo/infra"
)

func main() {
	storer := infra.NewFileStore("tasks.json")
	manager := todo.NewManager(storer)
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
			task, err := manager.AddTask(cmd.Add)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Printf("Added task %d: %s\n", task.Id, task.Title)
			}
		case cmd.List != nil:
			tasks, err := manager.ListTasks(cmd.List)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				printTasks(tasks)
			}
		case cmd.Delete != nil:
			if err := manager.DeleteTask(cmd.Delete); err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Printf("Deleted task #%d\n", cmd.Delete.Id)
			}
		case cmd.Update != nil:
			task, err := manager.UpdateTasks(cmd.Update)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Printf("Updated task #%d\n", task.Id)
			}
		case cmd.Done != nil:
			task, err := manager.ToggleDone(cmd.Done)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				status := "pending"
				if task.Done {
					status = "done"
				}
				fmt.Printf("Mark task #%d as %s\n", task.Id, status)
			}
		}
	}
}

func printTasks(tasks []todo.Task) {
	if len(tasks) == 0 {
		fmt.Println("No tasks found")
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tStatus\tPriority\tTitle\tDescription")
	fmt.Fprintln(w, "--\t------\t--------\t-----\t-----------")
	for _, t := range tasks {
		status := "[ ]"
		if t.Done {
			status = "[x]"
		}
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n", t.Id, status, t.Priority, t.Title, t.Description)
	}
	w.Flush()
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
