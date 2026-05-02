package todo

import (
	"fmt"
	"time"
)

func (t *Todo) AddTask(cmd *AddCmd) (*Task, error) {
	priority := cmd.Priority
	if priority == "" {
		priority = "low"
	}

	task := Task{
		Title:       cmd.Title,
		Description: cmd.Description,
		Priority:    priority,
		Done:        false,
		CreatedAt:   time.Now(),
	}

	created, err := t.store.Create(task)
	if err != nil {
		return nil, fmt.Errorf("create task: %w", err)
	}

	return &created, nil
}

func (t *Todo) ListTasks(cmd *ListCmd) ([]Task, error) {
	tasks, err := t.store.List(cmd.Filter)
	if err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}
	return tasks, nil
}

func (t *Todo) UpdateTasks(cmd *UpdateCmd) (*Task, error) {
	task := Task{
		Id:          cmd.Id,
		Title:       cmd.Title,
		Description: cmd.Description,
		Priority:    cmd.Priority,
	}

	updated, err := t.store.Update(task)
	if err != nil {
		return nil, fmt.Errorf("update task: %w", err)
	}

	return &updated, nil
}

func (t *Todo) DeleteTask(cmd *DeleteCmd) error {
	if err := t.store.Delete(cmd.Id); err != nil {
		return fmt.Errorf("delete task: %w", err)
	}
	return nil
}

func (t *Todo) ToggleDone(cmd *DoneCmd) (*Task, error) {
	task, err := t.store.ToggleDone(cmd.Id)
	if err != nil {
		return nil, fmt.Errorf("toggle done: %w", err)
	}
	return &task, nil
}
