package todo

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"
)

func nextId(tasks []Task) int {
	maxID := 0
	for _, t := range tasks {
		if t.Id > maxID {
			maxID = t.Id
		}
	}
	return maxID + 1
}

func (t *Todo) AddTask(cmd *AddCmd) error {
	tasks, err := t.store.Load()
	if err != nil {
		return fmt.Errorf("load tasks: %w", err)
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
	if err := t.store.Save(tasks); err != nil {
		return fmt.Errorf("save tasks: %w", err)
	}

	fmt.Printf("Added task %d: %s\n", task.Id, task.Title)
	return nil
}

func (t *Todo) ListTasks(cmd *ListCmd) error {
	tasks, err := t.store.Load()
	if err != nil {
		return fmt.Errorf("load tasks: %w", err)
	}

	if len(tasks) == 0 {
		fmt.Println("No tasks found")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tStatus\tPriority\tTitle\tDescription")
	fmt.Fprintln(w, "--\t------\t--------\t-----\t-----------")

	for _, task := range tasks {
		if cmd.Filter == "done" && !task.Done {
			continue
		}
		if cmd.Filter == "pending" && task.Done {
			continue
		}

		status := "[ ]"
		if task.Done {
			status = "[x]"
		}

		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n", task.Id, status, task.Priority, task.Title, task.Description)
	}
	w.Flush()
	return nil
}

func (t *Todo) UpdateTasks(cmd *UpdateCmd) error {
	tasks, err := t.store.Load()
	if err != nil {
		return fmt.Errorf("load tasks: %w", err)
	}

	for i, task := range tasks {
		if task.Id == cmd.Id {
			if cmd.Title != "" {
				tasks[i].Title = cmd.Title
			}
			if cmd.Description != "" {
				tasks[i].Description = cmd.Description
			}
			if cmd.Priority != "" {
				tasks[i].Priority = cmd.Priority
			}

			if err := t.store.Save(tasks); err != nil {
				return fmt.Errorf("save tasks: %w", err)
			}

			fmt.Printf("update task #%d successfully\n", task.Id)
			return nil
		}
	}
	return fmt.Errorf("task #%d not found", cmd.Id)
}

func (t *Todo) DeleteTask(cmd *DeleteCmd) error {
	tasks, err := t.store.Load()
	if err != nil {
		return fmt.Errorf("load tasks: %w", err)
	}

	for i, task := range tasks {
		if task.Id == cmd.Id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			if err := t.store.Save(tasks); err != nil {
				return fmt.Errorf("save tasks: %w", err)
			}

			fmt.Printf("Deleted task #%d\n", task.Id)
			return nil
		}
	}
	return fmt.Errorf("task #%d not found", cmd.Id)
}

func (t *Todo) ToggleDone(cmd *DoneCmd) error {
	tasks, err := t.store.Load()
	if err != nil {
		return fmt.Errorf("load tasks: %w", err)
	}

	for i, task := range tasks {
		if task.Id == cmd.Id {
			tasks[i].Done = !tasks[i].Done
			if err := t.store.Save(tasks); err != nil {
				return fmt.Errorf("save tasks: %w", err)
			}
			status := "pending"
			if tasks[i].Done {
				status = "done"
			}
			fmt.Printf("Mark task #%d as %s\n", tasks[i].Id, status)
			return nil
		}
	}
	return fmt.Errorf("task #%d not found", cmd.Id)
}
