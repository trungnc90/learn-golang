package infra

import (
	"os"
	"testing"

	todo "github.com/trungnc90/learn-golang/assignments/01_todo"
)

func newTestFileStore(t *testing.T) *FileStore {
	t.Helper()
	path := "test_tasks.json"
	os.Remove(path)
	t.Cleanup(func() { os.Remove(path) })
	return NewFileStore(path)
}

func TestFileStore(t *testing.T) {
	t.Run("create assigns id", func(t *testing.T) {
		store := newTestFileStore(t)
		task, err := store.Create(todo.Task{Title: "Task A", Priority: "high"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if task.Id != 1 {
			t.Fatalf("expected id 1, got %d", task.Id)
		}
		if task.Title != "Task A" {
			t.Fatalf("expected 'Task A', got '%s'", task.Title)
		}
	})

	t.Run("create auto-increments id", func(t *testing.T) {
		store := newTestFileStore(t)
		store.Create(todo.Task{Title: "First"})
		task, _ := store.Create(todo.Task{Title: "Second"})
		if task.Id != 2 {
			t.Fatalf("expected id 2, got %d", task.Id)
		}
	})

	t.Run("list returns all tasks", func(t *testing.T) {
		store := newTestFileStore(t)
		store.Create(todo.Task{Title: "Task A"})
		store.Create(todo.Task{Title: "Task B"})

		tasks, err := store.List("")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(tasks) != 2 {
			t.Fatalf("expected 2 tasks, got %d", len(tasks))
		}
	})

	t.Run("list filters done tasks", func(t *testing.T) {
		store := newTestFileStore(t)
		store.Create(todo.Task{Title: "Pending", Done: false})
		store.Create(todo.Task{Title: "Done", Done: true})

		tasks, _ := store.List("done")
		if len(tasks) != 1 || tasks[0].Title != "Done" {
			t.Fatalf("expected only done task")
		}
	})

	t.Run("list from missing file returns empty", func(t *testing.T) {
		store := newTestFileStore(t)
		tasks, err := store.List("")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(tasks) != 0 {
			t.Fatalf("expected 0 tasks, got %d", len(tasks))
		}
	})

	t.Run("get by id returns task", func(t *testing.T) {
		store := newTestFileStore(t)
		store.Create(todo.Task{Title: "Task A"})

		task, err := store.GetByID(1)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if task.Title != "Task A" {
			t.Fatalf("expected 'Task A', got '%s'", task.Title)
		}
	})

	t.Run("get by id returns error when not found", func(t *testing.T) {
		store := newTestFileStore(t)
		_, err := store.GetByID(999)
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("update modifies task", func(t *testing.T) {
		store := newTestFileStore(t)
		store.Create(todo.Task{Title: "Old", Priority: "low"})

		updated, err := store.Update(todo.Task{Id: 1, Title: "New", Priority: "high"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if updated.Title != "New" {
			t.Fatalf("expected 'New', got '%s'", updated.Title)
		}
		if updated.Priority != "high" {
			t.Fatalf("expected 'high', got '%s'", updated.Priority)
		}
	})

	t.Run("update returns error when not found", func(t *testing.T) {
		store := newTestFileStore(t)
		_, err := store.Update(todo.Task{Id: 999, Title: "Nope"})
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("delete removes task", func(t *testing.T) {
		store := newTestFileStore(t)
		store.Create(todo.Task{Title: "Task A"})
		store.Create(todo.Task{Title: "Task B"})

		err := store.Delete(1)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		tasks, _ := store.List("")
		if len(tasks) != 1 {
			t.Fatalf("expected 1 task, got %d", len(tasks))
		}
	})

	t.Run("delete returns error when not found", func(t *testing.T) {
		store := newTestFileStore(t)
		err := store.Delete(999)
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("toggle done flips status", func(t *testing.T) {
		store := newTestFileStore(t)
		store.Create(todo.Task{Title: "Task", Done: false})

		task, err := store.ToggleDone(1)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !task.Done {
			t.Fatal("expected done=true")
		}
	})

	t.Run("toggle done returns error when not found", func(t *testing.T) {
		store := newTestFileStore(t)
		_, err := store.ToggleDone(999)
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("load corrupted file returns error", func(t *testing.T) {
		path := "test_corrupt.json"
		os.WriteFile(path, []byte("not valid json"), 0644)
		t.Cleanup(func() { os.Remove(path) })
		store := NewFileStore(path)
		_, err := store.List("")
		if err == nil {
			t.Fatal("expected error for corrupted JSON file")
		}
	})
}
