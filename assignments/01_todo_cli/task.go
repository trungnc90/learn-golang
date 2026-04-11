package main

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
	"time"
)

const dataFile = "tasks.json"

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

func saveTasks(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(dataFile, data, 0644)
}

func loadTasks() ([]Task, error) {
	data, err := os.ReadFile(dataFile)

	if err != nil {
		if os.IsNotExist(err) {
			return []Task{}, nil
		}
		return nil, err
	}

	var tasks []Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

func addTask(title, description, priority string) {
	tasks, err := loadTasks()
	if err != nil {
		fmt.Println("Error loading current tasks:", err)
		return
	}

	if priority == "" {
		priority = "low"
	}

	task := Task{
		Id:          nextId(tasks),
		Title:       title,
		Description: description,
		Priority:    priority,
		Done:        false,
		CreatedAt:   time.Now(),
	}

	tasks = append(tasks, task)
	if err := saveTasks(tasks); err != nil {
		fmt.Println("Error saving tasks:", err)
		return
	}

	fmt.Printf("Added task %d: %s\n", task.Id, task.Title)
}

func listTasks(filter string) {
	tasks, err := loadTasks()
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
		if filter == "done" && !t.Done {
			//skip not done tasks
			continue
		}
		if filter == "pending" && t.Done {
			// skip done tasks
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

func updateTasks(id int, title string, description string, priority string) {
	tasks, err := loadTasks()
	if err != nil {
		fmt.Println("updateTasks(): Error loading tasks:", err)
		return
	}

	for i, t := range tasks {
		if t.Id == id {
			if title != "" {
				tasks[i].Title = title
			}
			if description != "" {
				tasks[i].Description = description
			}
			if priority != "" {
				tasks[i].Priority = priority
			}

			if err := saveTasks(tasks); err != nil {
				fmt.Println("updateTasks(): Error saving tasks", err)
				return
			}

			fmt.Printf("update task #%d successfully\n", t.Id)
			return
		}
	}
	fmt.Printf("Task #%d not found\n", id)
}

func deleteTask(id int) {
	tasks, err := loadTasks()
	if err != nil {
		fmt.Println("deleteTask(): Error loading tasks:", err)
		return
	}

	for i, t := range tasks {
		if t.Id == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			if err := saveTasks(tasks); err != nil {
				fmt.Println("Error saving tasks:", err)
			}

			fmt.Printf("Deleted task #%d\n", t.Id)
			return
		}
	}
	fmt.Printf("Task #%d not found\n", id)
}

func toggleDone(id int) {
	tasks, err := loadTasks()
	if err != nil {
		fmt.Println("toggleDone(): Error loading tasks", err)
		return
	}

	for i, t := range tasks {
		if t.Id == id {
			tasks[i].Done = !tasks[i].Done
			if err := saveTasks(tasks); err != nil {
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
	fmt.Printf("Task #%d not found\n", id)
}
