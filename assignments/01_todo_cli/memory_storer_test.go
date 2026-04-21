package todo

import (
	"testing"
)

func TestMemoryStore(t *testing.T) {
	t.Run("load from empty store", func(t *testing.T) {
		store := NewMemoryStore()
		tasks, err := store.Load()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(tasks) != 0 {
			t.Fatalf("expected 0 tasks, got %d", len(tasks))
		}
	})

	t.Run("save and load round trip", func(t *testing.T) {
		store := NewMemoryStore()
		original := []Task{
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
		store := NewMemoryStore()
		store.Save([]Task{{Id: 1, Title: "Old"}})
		store.Save([]Task{{Id: 1, Title: "New"}, {Id: 2, Title: "Extra"}})
		tasks, _ := store.Load()
		if len(tasks) != 2 {
			t.Fatalf("expected 2 tasks, got %d", len(tasks))
		}
		if tasks[0].Title != "New" {
			t.Fatalf("expected 'New', got '%s'", tasks[0].Title)
		}
	})

	t.Run("save empty clears data", func(t *testing.T) {
		store := NewMemoryStore()
		store.Save([]Task{{Id: 1, Title: "Task"}})
		store.Save([]Task{})
		tasks, _ := store.Load()
		if len(tasks) != 0 {
			t.Fatalf("expected 0 tasks, got %d", len(tasks))
		}
	})

	t.Run("load returns a copy", func(t *testing.T) {
		store := NewMemoryStore()
		store.Save([]Task{{Id: 1, Title: "Original"}})
		loaded, _ := store.Load()
		loaded[0].Title = "Mutated"
		fresh, _ := store.Load()
		if fresh[0].Title != "Original" {
			t.Fatalf("expected 'Original', got '%s' — Load did not return a copy", fresh[0].Title)
		}
	})

	t.Run("save stores a copy", func(t *testing.T) {
		store := NewMemoryStore()
		tasks := []Task{{Id: 1, Title: "Original"}}
		store.Save(tasks)
		tasks[0].Title = "Mutated"
		loaded, _ := store.Load()
		if loaded[0].Title != "Original" {
			t.Fatalf("expected 'Original', got '%s' — Save did not copy input", loaded[0].Title)
		}
	})

	t.Run("save never returns error", func(t *testing.T) {
		store := NewMemoryStore()
		err := store.Save([]Task{{Id: 1}})
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
	})

	t.Run("load never returns error", func(t *testing.T) {
		store := NewMemoryStore()
		_, err := store.Load()
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
	})
}
