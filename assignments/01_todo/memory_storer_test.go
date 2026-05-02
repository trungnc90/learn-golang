package todo

import (
	"testing"
)

func TestMemoryStore(t *testing.T) {
	t.Run("create assigns id", func(t *testing.T) {
		store := NewMemoryStore()
		task, err := store.Create(Task{Title: "Task A", Priority: "high"})
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
		store := NewMemoryStore()
		store.Create(Task{Title: "First"})
		task, _ := store.Create(Task{Title: "Second"})
		if task.Id != 2 {
			t.Fatalf("expected id 2, got %d", task.Id)
		}
	})

	t.Run("list returns all tasks", func(t *testing.T) {
		store := NewMemoryStore()
		store.Create(Task{Title: "Task A"})
		store.Create(Task{Title: "Task B"})

		tasks, err := store.List("")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(tasks) != 2 {
			t.Fatalf("expected 2 tasks, got %d", len(tasks))
		}
	})

	t.Run("list filters done tasks", func(t *testing.T) {
		store := NewMemoryStore()
		store.Create(Task{Title: "Pending", Done: false})
		store.Create(Task{Title: "Done", Done: true})

		tasks, _ := store.List("done")
		if len(tasks) != 1 || tasks[0].Title != "Done" {
			t.Fatalf("expected only done task")
		}
	})

	t.Run("list filters pending tasks", func(t *testing.T) {
		store := NewMemoryStore()
		store.Create(Task{Title: "Pending", Done: false})
		store.Create(Task{Title: "Done", Done: true})

		tasks, _ := store.List("pending")
		if len(tasks) != 1 || tasks[0].Title != "Pending" {
			t.Fatalf("expected only pending task")
		}
	})

	t.Run("list empty store returns empty slice", func(t *testing.T) {
		store := NewMemoryStore()
		tasks, err := store.List("")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if tasks == nil {
			t.Fatal("expected empty slice, got nil")
		}
		if len(tasks) != 0 {
			t.Fatalf("expected 0 tasks, got %d", len(tasks))
		}
	})

	t.Run("get by id returns task", func(t *testing.T) {
		store := NewMemoryStore()
		store.Create(Task{Title: "Task A"})

		task, err := store.GetByID(1)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if task.Title != "Task A" {
			t.Fatalf("expected 'Task A', got '%s'", task.Title)
		}
	})

	t.Run("get by id returns error when not found", func(t *testing.T) {
		store := NewMemoryStore()
		_, err := store.GetByID(999)
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("update modifies task", func(t *testing.T) {
		store := NewMemoryStore()
		store.Create(Task{Title: "Old", Priority: "low"})

		updated, err := store.Update(Task{Id: 1, Title: "New", Priority: "high"})
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
		store := NewMemoryStore()
		_, err := store.Update(Task{Id: 999, Title: "Nope"})
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("delete removes task", func(t *testing.T) {
		store := NewMemoryStore()
		store.Create(Task{Title: "Task A"})
		store.Create(Task{Title: "Task B"})

		err := store.Delete(1)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		tasks, _ := store.List("")
		if len(tasks) != 1 {
			t.Fatalf("expected 1 task, got %d", len(tasks))
		}
		if tasks[0].Title != "Task B" {
			t.Fatalf("expected 'Task B', got '%s'", tasks[0].Title)
		}
	})

	t.Run("delete returns error when not found", func(t *testing.T) {
		store := NewMemoryStore()
		err := store.Delete(999)
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("toggle done flips status", func(t *testing.T) {
		store := NewMemoryStore()
		store.Create(Task{Title: "Task", Done: false})

		task, err := store.ToggleDone(1)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !task.Done {
			t.Fatal("expected done=true")
		}

		task, _ = store.ToggleDone(1)
		if task.Done {
			t.Fatal("expected done=false after second toggle")
		}
	})

	t.Run("toggle done returns error when not found", func(t *testing.T) {
		store := NewMemoryStore()
		_, err := store.ToggleDone(999)
		if err == nil {
			t.Fatal("expected error")
		}
	})
}
