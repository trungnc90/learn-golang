// How to run tests:
//   go test -v ./...                    Run all tests
//   go test -v -run TestAddTask ./...   Run addTask subtests only
//   go test -count=1 ./...              Run without cache

package todo

import (
	"fmt"
	"testing"
)

// newTestTodo creates a Todo with a mock storer for testing.
func newTestTodo() (*Todo, *storer) {
	mock := testStorer()
	return NewManager(mock), mock
}

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

// --- AddTask ---

func TestAddTask(t *testing.T) {
	t.Run("returns created task", func(t *testing.T) {
		app, mock := newTestTodo()
		mock.EXPECT().Load(t).Return([]Task{}, nil)
		mock.EXPECT().Save(t).MatchTasks(func(tasks []Task) bool {
			return len(tasks) == 1 && tasks[0].Title == "Buy groceries"
		}).Return(nil)

		task, err := app.AddTask(&AddCmd{Title: "Buy groceries", Description: "Milk and eggs", Priority: "high"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if task.Title != "Buy groceries" {
			t.Fatalf("got '%s'", task.Title)
		}
		if task.Priority != "high" {
			t.Fatalf("got '%s'", task.Priority)
		}
		if task.Id != 1 {
			t.Fatalf("expected id 1, got %d", task.Id)
		}
	})

	t.Run("default priority is low", func(t *testing.T) {
		app, mock := newTestTodo()
		mock.EXPECT().Load(t).Return([]Task{}, nil)
		mock.EXPECT().Save(t).Return(nil)

		task, err := app.AddTask(&AddCmd{Title: "Task"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if task.Priority != "low" {
			t.Fatalf("got '%s'", task.Priority)
		}
	})

	t.Run("returns error on load failure", func(t *testing.T) {
		app, mock := newTestTodo()
		mock.EXPECT().Load(t).Return(nil, fmt.Errorf("disk error"))

		_, err := app.AddTask(&AddCmd{Title: "Test"})
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("returns error on save failure", func(t *testing.T) {
		app, mock := newTestTodo()
		mock.EXPECT().Load(t).Return([]Task{}, nil)
		mock.EXPECT().Save(t).Return(fmt.Errorf("write failed"))

		_, err := app.AddTask(&AddCmd{Title: "Test"})
		if err == nil {
			t.Fatal("expected error")
		}
	})
}

// --- ListTasks ---

func TestListTasks(t *testing.T) {
	t.Run("returns all tasks", func(t *testing.T) {
		app, mock := newTestTodo()
		mock.EXPECT().Load(t).Return([]Task{
			{Id: 1, Title: "Task 1"},
			{Id: 2, Title: "Task 2", Done: true},
		}, nil)

		tasks, err := app.ListTasks(&ListCmd{})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(tasks) != 2 {
			t.Fatalf("expected 2, got %d", len(tasks))
		}
	})

	t.Run("filters done tasks", func(t *testing.T) {
		app, mock := newTestTodo()
		mock.EXPECT().Load(t).Return([]Task{
			{Id: 1, Title: "Pending", Done: false},
			{Id: 2, Title: "Done", Done: true},
		}, nil)

		tasks, err := app.ListTasks(&ListCmd{Filter: "done"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(tasks) != 1 || tasks[0].Title != "Done" {
			t.Fatalf("expected only done task")
		}
	})

	t.Run("filters pending tasks", func(t *testing.T) {
		app, mock := newTestTodo()
		mock.EXPECT().Load(t).Return([]Task{
			{Id: 1, Title: "Pending", Done: false},
			{Id: 2, Title: "Done", Done: true},
		}, nil)

		tasks, err := app.ListTasks(&ListCmd{Filter: "pending"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(tasks) != 1 || tasks[0].Title != "Pending" {
			t.Fatalf("expected only pending task")
		}
	})

	t.Run("returns error on load failure", func(t *testing.T) {
		app, mock := newTestTodo()
		mock.EXPECT().Load(t).Return(nil, fmt.Errorf("error"))

		_, err := app.ListTasks(&ListCmd{})
		if err == nil {
			t.Fatal("expected error")
		}
	})
}

// --- ToggleDone ---

func TestToggleDone(t *testing.T) {
	t.Run("returns toggled task", func(t *testing.T) {
		app, mock := newTestTodo()
		mock.EXPECT().Load(t).Return([]Task{{Id: 1, Done: false}}, nil)
		mock.EXPECT().Save(t).Return(nil)

		task, err := app.ToggleDone(&DoneCmd{Id: 1})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !task.Done {
			t.Fatal("expected done")
		}
	})

	t.Run("returns error when not found", func(t *testing.T) {
		app, mock := newTestTodo()
		mock.EXPECT().Load(t).Return([]Task{}, nil)

		_, err := app.ToggleDone(&DoneCmd{Id: 999})
		if err == nil {
			t.Fatal("expected error")
		}
	})
}

// --- DeleteTask ---

func TestDeleteTask(t *testing.T) {
	t.Run("removes task by id", func(t *testing.T) {
		app, mock := newTestTodo()
		mock.EXPECT().Load(t).Return([]Task{{Id: 1}, {Id: 2}}, nil)
		mock.EXPECT().Save(t).MatchTasks(func(tasks []Task) bool {
			return len(tasks) == 1 && tasks[0].Id == 2
		}).Return(nil)

		err := app.DeleteTask(&DeleteCmd{Id: 1})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("returns error when not found", func(t *testing.T) {
		app, mock := newTestTodo()
		mock.EXPECT().Load(t).Return([]Task{}, nil)

		err := app.DeleteTask(&DeleteCmd{Id: 999})
		if err == nil {
			t.Fatal("expected error")
		}
	})
}

// --- UpdateTasks ---

func TestUpdateTasks(t *testing.T) {
	t.Run("returns updated task", func(t *testing.T) {
		app, mock := newTestTodo()
		mock.EXPECT().Load(t).Return([]Task{{Id: 1, Title: "Old", Description: "desc", Priority: "low"}}, nil)
		mock.EXPECT().Save(t).Return(nil)

		task, err := app.UpdateTasks(&UpdateCmd{Id: 1, Title: "New", Priority: "high"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if task.Title != "New" {
			t.Fatalf("got '%s'", task.Title)
		}
		if task.Description != "desc" {
			t.Fatalf("description changed: '%s'", task.Description)
		}
		if task.Priority != "high" {
			t.Fatalf("got '%s'", task.Priority)
		}
	})

	t.Run("returns error when not found", func(t *testing.T) {
		app, mock := newTestTodo()
		mock.EXPECT().Load(t).Return([]Task{}, nil)

		_, err := app.UpdateTasks(&UpdateCmd{Id: 999, Title: "title"})
		if err == nil {
			t.Fatal("expected error")
		}
	})
}
