// How to run tests:
//   go test -v ./...                    Run all tests
//   go test -v -run TestAddTask ./...   Run addTask subtests only
//   go test -count=1 ./...              Run without cache

package todo

import (
	"testing"
)

// --- nextId ---

func TestNextId(t *testing.T) {
	t.Run("empty list returns 1", func(t *testing.T) {
		id := nextId([]Task{})
		if id != 1 {
			t.Fatalf("expected 1, got %d", id)
		}
	})

	t.Run("returns max id + 1", func(t *testing.T) {
		tasks := []Task{{Id: 1}, {Id: 5}, {Id: 3}}
		id := nextId(tasks)
		if id != 6 {
			t.Fatalf("expected 6, got %d", id)
		}
	})
}

func TestAddTask(t *testing.T) {
	t.Run("basic add with all fields", func(t *testing.T) {
		store := NewMemoryStore()
		AddTask(store, &AddCmd{Title: "Buy groceries", Description: "Milk and eggs", Priority: "high"})
		tasks, _ := store.Load()
		if len(tasks) != 1 {
			t.Fatalf("expected 1 task, got %d", len(tasks))
		}
		if tasks[0].Title != "Buy groceries" {
			t.Fatalf("got '%s'", tasks[0].Title)
		}
		if tasks[0].Description != "Milk and eggs" {
			t.Fatalf("got '%s'", tasks[0].Description)
		}
		if tasks[0].Priority != "high" {
			t.Fatalf("got '%s'", tasks[0].Priority)
		}
		if tasks[0].Done != false {
			t.Fatal("expected not done")
		}
	})
	t.Run("default priority is low", func(t *testing.T) {
		store := NewMemoryStore()
		AddTask(store, &AddCmd{Title: "No priority task"})
		tasks, _ := store.Load()
		if tasks[0].Priority != "low" {
			t.Fatalf("got '%s'", tasks[0].Priority)
		}
	})
	t.Run("multiple tasks get sequential IDs", func(t *testing.T) {
		store := NewMemoryStore()
		AddTask(store, &AddCmd{Title: "Task 1", Priority: "high"})
		AddTask(store, &AddCmd{Title: "Task 2", Priority: "medium"})
		AddTask(store, &AddCmd{Title: "Task 3", Priority: "low"})
		tasks, _ := store.Load()
		if len(tasks) != 3 {
			t.Fatalf("expected 3, got %d", len(tasks))
		}
		if tasks[0].Id != 1 || tasks[1].Id != 2 || tasks[2].Id != 3 {
			t.Fatal("IDs not sequential")
		}
	})
}

func TestToggleDone(t *testing.T) {
	t.Run("toggle on and off", func(t *testing.T) {
		store := NewMemoryStore()
		AddTask(store, &AddCmd{Title: "Toggle me", Priority: "low"})
		ToggleDone(store, &DoneCmd{Id: 1})
		tasks, _ := store.Load()
		if !tasks[0].Done {
			t.Fatal("expected done")
		}
		ToggleDone(store, &DoneCmd{Id: 1})
		tasks, _ = store.Load()
		if tasks[0].Done {
			t.Fatal("expected pending")
		}
	})
	t.Run("not found", func(t *testing.T) {
		store := NewMemoryStore()
		ToggleDone(store, &DoneCmd{Id: 999})
	})
}

func TestDeleteTask(t *testing.T) {
	t.Run("delete by id", func(t *testing.T) {
		store := NewMemoryStore()
		AddTask(store, &AddCmd{Title: "Task 1", Priority: "high"})
		AddTask(store, &AddCmd{Title: "Task 2", Priority: "low"})
		DeleteTask(store, &DeleteCmd{Id: 1})
		tasks, _ := store.Load()
		if len(tasks) != 1 {
			t.Fatalf("expected 1, got %d", len(tasks))
		}
		if tasks[0].Title != "Task 2" {
			t.Fatalf("got '%s'", tasks[0].Title)
		}
	})
	t.Run("not found", func(t *testing.T) {
		store := NewMemoryStore()
		DeleteTask(store, &DeleteCmd{Id: 999})
	})
}

func TestUpdateTasks(t *testing.T) {
	t.Run("update title only", func(t *testing.T) {
		store := NewMemoryStore()
		AddTask(store, &AddCmd{Title: "Original", Description: "desc", Priority: "low"})
		UpdateTasks(store, &UpdateCmd{Id: 1, Title: "Updated Title"})
		tasks, _ := store.Load()
		if tasks[0].Title != "Updated Title" {
			t.Fatalf("got '%s'", tasks[0].Title)
		}
		if tasks[0].Description != "desc" {
			t.Fatalf("got '%s'", tasks[0].Description)
		}
		if tasks[0].Priority != "low" {
			t.Fatalf("got '%s'", tasks[0].Priority)
		}
	})
	t.Run("update all fields", func(t *testing.T) {
		store := NewMemoryStore()
		AddTask(store, &AddCmd{Title: "Original", Description: "old desc", Priority: "low"})
		UpdateTasks(store, &UpdateCmd{Id: 1, Title: "New Title", Description: "new desc", Priority: "high"})
		tasks, _ := store.Load()
		if tasks[0].Title != "New Title" {
			t.Fatalf("got '%s'", tasks[0].Title)
		}
		if tasks[0].Description != "new desc" {
			t.Fatalf("got '%s'", tasks[0].Description)
		}
		if tasks[0].Priority != "high" {
			t.Fatalf("got '%s'", tasks[0].Priority)
		}
	})
	t.Run("not found", func(t *testing.T) {
		store := NewMemoryStore()
		UpdateTasks(store, &UpdateCmd{Id: 999, Title: "title"})
	})
}

func TestListTasks(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		store := NewMemoryStore()
		ListTasks(store, &ListCmd{})
	})
	t.Run("with filters", func(t *testing.T) {
		store := NewMemoryStore()
		AddTask(store, &AddCmd{Title: "Pending task", Priority: "low"})
		AddTask(store, &AddCmd{Title: "Done task", Priority: "high"})
		ToggleDone(store, &DoneCmd{Id: 2})
		ListTasks(store, &ListCmd{})
		ListTasks(store, &ListCmd{Filter: "done"})
		ListTasks(store, &ListCmd{Filter: "pending"})
	})
}
