// How to run tests:
//   go test -v ./...          Run all tests with verbose output
//   go test -run TestAddTask  Run a specific test by name
//   go test -count=1 ./...    Run without cache

package main

import (
	"os"
	"testing"
)

// setup creates a clean test environment by removing any existing tasks.json
// and returns a cleanup function to call with defer.
func setup(t *testing.T) func() {
	t.Helper()
	os.Remove(dataFile)
	return func() {
		os.Remove(dataFile)
	}
}

// --- loadTasks / saveTasks ---

func TestLoadTasks_NoFile(t *testing.T) {
	cleanup := setup(t)
	defer cleanup()

	tasks, err := loadTasks()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(tasks) != 0 {
		t.Fatalf("expected 0 tasks, got %d", len(tasks))
	}
}

func TestSaveAndLoadTasks(t *testing.T) {
	cleanup := setup(t)
	defer cleanup()

	original := []Task{
		{Id: 1, Title: "Test", Priority: "high"},
		{Id: 2, Title: "Test2", Priority: "low"},
	}
	if err := saveTasks(original); err != nil {
		t.Fatalf("saveTasks failed: %v", err)
	}

	loaded, err := loadTasks()
	if err != nil {
		t.Fatalf("loadTasks failed: %v", err)
	}
	if len(loaded) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(loaded))
	}
	if loaded[0].Title != "Test" || loaded[1].Title != "Test2" {
		t.Fatalf("task titles don't match")
	}
}

// --- nextId ---

func TestNextId_Empty(t *testing.T) {
	id := nextId([]Task{})
	if id != 1 {
		t.Fatalf("expected 1, got %d", id)
	}
}

func TestNextId_WithTasks(t *testing.T) {
	tasks := []Task{
		{Id: 1}, {Id: 5}, {Id: 3},
	}
	id := nextId(tasks)
	if id != 6 {
		t.Fatalf("expected 6, got %d", id)
	}
}

// --- addTask ---

func TestAddTask_Basic(t *testing.T) {
	cleanup := setup(t)
	defer cleanup()

	addTask("Buy groceries", "Milk and eggs", "high")

	tasks, _ := loadTasks()
	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(tasks))
	}
	if tasks[0].Title != "Buy groceries" {
		t.Fatalf("expected title 'Buy groceries', got '%s'", tasks[0].Title)
	}
	if tasks[0].Description != "Milk and eggs" {
		t.Fatalf("expected description 'Milk and eggs', got '%s'", tasks[0].Description)
	}
	if tasks[0].Priority != "high" {
		t.Fatalf("expected priority 'high', got '%s'", tasks[0].Priority)
	}
	if tasks[0].Done != false {
		t.Fatal("expected task to not be done")
	}
}

func TestAddTask_DefaultPriority(t *testing.T) {
	cleanup := setup(t)
	defer cleanup()

	addTask("No priority task", "", "")

	tasks, _ := loadTasks()
	if tasks[0].Priority != "low" {
		t.Fatalf("expected default priority 'low', got '%s'", tasks[0].Priority)
	}
}

func TestAddTask_MultipleTasks(t *testing.T) {
	cleanup := setup(t)
	defer cleanup()

	addTask("Task 1", "", "high")
	addTask("Task 2", "", "medium")
	addTask("Task 3", "", "low")

	tasks, _ := loadTasks()
	if len(tasks) != 3 {
		t.Fatalf("expected 3 tasks, got %d", len(tasks))
	}
	if tasks[0].Id != 1 || tasks[1].Id != 2 || tasks[2].Id != 3 {
		t.Fatal("task IDs are not sequential")
	}
}

// --- toggleDone ---

func TestToggleDone(t *testing.T) {
	cleanup := setup(t)
	defer cleanup()

	addTask("Toggle me", "", "low")

	// Mark as done
	toggleDone(1)
	tasks, _ := loadTasks()
	if !tasks[0].Done {
		t.Fatal("expected task to be done")
	}

	// Toggle back to pending
	toggleDone(1)
	tasks, _ = loadTasks()
	if tasks[0].Done {
		t.Fatal("expected task to be pending")
	}
}

func TestToggleDone_NotFound(t *testing.T) {
	cleanup := setup(t)
	defer cleanup()

	// Should not panic on non-existent task
	toggleDone(999)
}

// --- deleteTask ---

func TestDeleteTask(t *testing.T) {
	cleanup := setup(t)
	defer cleanup()

	addTask("Task 1", "", "high")
	addTask("Task 2", "", "low")

	deleteTask(1)

	tasks, _ := loadTasks()
	if len(tasks) != 1 {
		t.Fatalf("expected 1 task after delete, got %d", len(tasks))
	}
	if tasks[0].Title != "Task 2" {
		t.Fatalf("expected remaining task 'Task 2', got '%s'", tasks[0].Title)
	}
}

func TestDeleteTask_NotFound(t *testing.T) {
	cleanup := setup(t)
	defer cleanup()

	// Should not panic on non-existent task
	deleteTask(999)
}

// --- updateTasks ---

func TestUpdateTasks_Title(t *testing.T) {
	cleanup := setup(t)
	defer cleanup()

	addTask("Original", "desc", "low")

	updateTasks(1, "Updated Title", "", "")

	tasks, _ := loadTasks()
	if tasks[0].Title != "Updated Title" {
		t.Fatalf("expected 'Updated Title', got '%s'", tasks[0].Title)
	}
	// Description and priority should remain unchanged
	if tasks[0].Description != "desc" {
		t.Fatalf("description should not change, got '%s'", tasks[0].Description)
	}
	if tasks[0].Priority != "low" {
		t.Fatalf("priority should not change, got '%s'", tasks[0].Priority)
	}
}

func TestUpdateTasks_AllFields(t *testing.T) {
	cleanup := setup(t)
	defer cleanup()

	addTask("Original", "old desc", "low")

	updateTasks(1, "New Title", "new desc", "high")

	tasks, _ := loadTasks()
	if tasks[0].Title != "New Title" {
		t.Fatalf("expected 'New Title', got '%s'", tasks[0].Title)
	}
	if tasks[0].Description != "new desc" {
		t.Fatalf("expected 'new desc', got '%s'", tasks[0].Description)
	}
	if tasks[0].Priority != "high" {
		t.Fatalf("expected 'high', got '%s'", tasks[0].Priority)
	}
}

func TestUpdateTasks_NotFound(t *testing.T) {
	cleanup := setup(t)
	defer cleanup()

	// Should not panic on non-existent task
	updateTasks(999, "title", "", "")
}

// --- listTasks ---

func TestListTasks_Empty(t *testing.T) {
	cleanup := setup(t)
	defer cleanup()

	// Should not panic with no tasks
	listTasks("")
}

func TestListTasks_WithFilter(t *testing.T) {
	cleanup := setup(t)
	defer cleanup()

	addTask("Pending task", "", "low")
	addTask("Done task", "", "high")
	toggleDone(2)

	// These should not panic — we can't easily capture stdout,
	// but we verify no crashes with filters
	listTasks("")
	listTasks("done")
	listTasks("pending")
}
