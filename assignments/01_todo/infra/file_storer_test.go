package infra

import (
	"os"
	"testing"

	todo "github.com/trungnc90/learn-golang/assignments/01_todo_cli"
)

func newTestFileStore(t *testing.T) *FileStore {
	t.Helper()
	path := "test_tasks.json"
	os.Remove(path)
	t.Cleanup(func() { os.Remove(path) })
	return NewFileStore(path)
}

func TestFileStore(t *testing.T) {
	t.Run("load from missing file", func(t *testing.T) {
		store := newTestFileStore(t)
		tasks, err := store.Load()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(tasks) != 0 {
			t.Fatalf("expected 0 tasks, got %d", len(tasks))
		}
	})

	t.Run("save and load round trip", func(t *testing.T) {
		store := newTestFileStore(t)
		original := []todo.Task{
			{Id: 1, Title: "Task A", Priority: "high"},
			{Id: 2, Title: "Task B", Priority: "low"},
		}
		if err := store.Save(original); err != nil {
			t.Fatalf("Save failed: %v", err)
		}
		loaded, err := store.Load()
		if err != nil {
			t.Fatalf("Load failed: %v", err)
		}
		if len(loaded) != 2 {
			t.Fatalf("expected 2 tasks, got %d", len(loaded))
		}
		if loaded[0].Title != "Task A" || loaded[1].Title != "Task B" {
			t.Fatal("task titles don't match")
		}
	})

	t.Run("save overwrites previous data", func(t *testing.T) {
		store := newTestFileStore(t)
		store.Save([]todo.Task{{Id: 1, Title: "Old"}})
		store.Save([]todo.Task{{Id: 1, Title: "New"}, {Id: 2, Title: "Extra"}})
		tasks, _ := store.Load()
		if len(tasks) != 2 {
			t.Fatalf("expected 2 tasks, got %d", len(tasks))
		}
		if tasks[0].Title != "New" {
			t.Fatalf("expected 'New', got '%s'", tasks[0].Title)
		}
	})

	t.Run("save empty clears data", func(t *testing.T) {
		store := newTestFileStore(t)
		store.Save([]todo.Task{{Id: 1, Title: "Task"}})
		store.Save([]todo.Task{})
		tasks, _ := store.Load()
		if len(tasks) != 0 {
			t.Fatalf("expected 0 tasks, got %d", len(tasks))
		}
	})

	t.Run("load returns a copy", func(t *testing.T) {
		store := newTestFileStore(t)
		store.Save([]todo.Task{{Id: 1, Title: "Original"}})
		loaded, _ := store.Load()
		loaded[0].Title = "Mutated"
		fresh, _ := store.Load()
		if fresh[0].Title != "Original" {
			t.Fatalf("expected 'Original', got '%s'", fresh[0].Title)
		}
	})

	t.Run("load corrupted file returns error", func(t *testing.T) {
		path := "test_corrupt.json"
		os.WriteFile(path, []byte("not valid json"), 0644)
		t.Cleanup(func() { os.Remove(path) })
		store := NewFileStore(path)
		_, err := store.Load()
		if err == nil {
			t.Fatal("expected error for corrupted JSON file")
		}
	})
}
