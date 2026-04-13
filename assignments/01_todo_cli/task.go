package main

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"
)

type Task struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Priority    string    `json:"priority"`
	Done        bool      `json:"done"`
	CreatedAt   time.Time `json:"created_at"`
}

func nextId(tasks []Task) int {
	maxID := 0
	for _, t := range tasks {
		if t.Id > maxID {
			maxID = t.Id
		}
	}
	return maxID + 1
}

func addTask(store Store, cmd *AddCmd) {
	tasks, err := store.Load()
	if err != nil {
		fmt.Println("Error loading current tasks:", err)
		return
	}

	priority := cmd.Priority
	if priority == "" {
		priority = "low"
	}

	task := Task{
		Id:          nextId(tasks),
		Title:       cmd.Title,
		Description: cmd.Description,
		Priority:    priority,
		Done:        false,
		CreatedAt:   time.Now(),
	}

	tasks = append(tasks, task)
	if err := store.Save(tasks); err != nil {
		fmt.Println("Error saving tasks:", err)
		return
	}

	fmt.Printf("Added task %d: %s\n", task.Id, task.Title)
}

func listTasks(store Store, cmd *ListCmd) {
	tasks, err := store.Load()
	if err != nil {
		fmt.Println("listTasks(): Error loading tasks:", err)
		return
	}

	if len(tasks) == 0 {
		fmt.Println("No tasks found")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tStatus\tPriority\tTitle\tDescription")
	fmt.Fprintln(w, "--\t------\t--------\t-----\t-----------")

	for _, t := range tasks {
		if cmd.Filter == "done" && !t.Done {
			continue
		}
		if cmd.Filter == "pending" && t.Done {
			continue
		}

		status := "[ ]"
		if t.Done {
			status = "[x]"
		}

		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n", t.Id, status, t.Priority, t.Title, t.Description)
	}
	w.Flush()
}

func updateTasks(store Store, cmd *UpdateCmd) {
	tasks, err := store.Load()
	if err != nil {
		fmt.Println("updateTasks(): Error loading tasks:", err)
		return
	}

	for i, t := range tasks {
		if t.Id == cmd.Id {
			if cmd.Title != "" {
				tasks[i].Title = cmd.Title
			}
			if cmd.Description != "" {
				tasks[i].Description = cmd.Description
			}
			if cmd.Priority != "" {
				tasks[i].Priority = cmd.Priority
			}

			if err := store.Save(tasks); err != nil {
				fmt.Println("updateTasks(): Error saving tasks", err)
				return
			}

			fmt.Printf("update task #%d successfully\n", t.Id)
			return
		}
	}
	fmt.Printf("Task #%d not found\n", cmd.Id)
}

func deleteTask(store Store, cmd *DeleteCmd) {
	tasks, err := store.Load()
	if err != nil {
		fmt.Println("deleteTask(): Error loading tasks:", err)
		return
	}

	for i, t := range tasks {
		if t.Id == cmd.Id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			if err := store.Save(tasks); err != nil {
				fmt.Println("Error saving tasks:", err)
			}

			fmt.Printf("Deleted task #%d\n", t.Id)
			return
		}
	}
	fmt.Printf("Task #%d not found\n", cmd.Id)
}

func toggleDone(store Store, cmd *DoneCmd) {
	tasks, err := store.Load()
	if err != nil {
		fmt.Println("toggleDone(): Error loading tasks", err)
		return
	}

	for i, t := range tasks {
		if t.Id == cmd.Id {
			tasks[i].Done = !tasks[i].Done
			if err := store.Save(tasks); err != nil {
				fmt.Println("toggleDone(): Error saving tasks", err)
				return
			}
			status := "pending"
			if tasks[i].Done {
				status = "done"
			}
			fmt.Printf("Mark task #%d as %s\n", tasks[i].Id, status)
			return
		}
	}
	fmt.Printf("Task #%d not found\n", cmd.Id)
}
