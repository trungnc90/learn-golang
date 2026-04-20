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
	return &Todo{store: mock}, mock
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
	t.Run("calls Load then Save with new task", func(t *testing.T) {
		app, mock := newTestTodo()
		mock.EXPECT().Load(t).Return([]Task{}, nil)
		mock.EXPECT().Save(t).MatchTasks(func(tasks []Task) bool {
			return len(tasks) == 1 &&
				tasks[0].Title == "Buy groceries" &&
				tasks[0].Description == "Milk and eggs" &&
				tasks[0].Priority == "high" &&
				tasks[0].Id == 1 &&
				!tasks[0].Done
		}).Return(nil)

		err := app.AddTask(&AddCmd{Title: "Buy groceries", Description: "Milk and eggs", Priority: "high"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("default priority is low", func(t *testing.T) {
		app, mock := newTestTodo()
		mock.EXPECT().Load(t).Return([]Task{}, nil)
		mock.EXPECT().Save(t).MatchTasks(func(tasks []Task) bool {
			return tasks[0].Priority == "low"
		}).Return(nil)

		err := app.AddTask(&AddCmd{Title: "Task"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("appends to existing tasks", func(t *testing.T) {
		app, mock := newTestTodo()
		mock.EXPECT().Load(t).Return([]Task{{Id: 3, Title: "Existing"}}, nil)
		mock.EXPECT().Save(t).MatchTasks(func(tasks []Task) bool {
			return len(tasks) == 2 && tasks[1].Id == 4
		}).Return(nil)

		err := app.AddTask(&AddCmd{Title: "New", Priority: "high"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("returns error on load failure", func(t *testing.T) {
		app, mock := newTestTodo()
		mock.EXPECT().Load(t).Return(nil, fmt.Errorf("disk error"))

		err := app.AddTask(&AddCmd{Title: "Test"})
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("returns error on save failure", func(t *testing.T) {
		app, mock := newTestTodo()
		mock.EXPECT().Load(t).Return([]Task{}, nil)
		mock.EXPECT().Save(t).Return(fmt.Errorf("write failed"))

		err := app.AddTask(&AddCmd{Title: "Test"})
		if err == nil {
			t.Fatal("expected error")
		}
	})
}

// --- ToggleDone ---

func TestToggleDone(t *testing.T) {
	t.Run("marks task as done", func(t *testing.T) {
		app, mock := newTestTodo()
		mock.EXPECT().Load(t).Return([]Task{{Id: 1, Title: "Task", Done: false}}, nil)
		mock.EXPECT().Save(t).MatchTasks(func(tasks []Task) bool {
			return tasks[0].Done == true
		}).Return(nil)

		err := app.ToggleDone(&DoneCmd{Id: 1})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("marks task as pending", func(t *testing.T) {
		app, mock := newTestTodo()
		mock.EXPECT().Load(t).Return([]Task{{Id: 1, Title: "Task", Done: true}}, nil)
		mock.EXPECT().Save(t).MatchTasks(func(tasks []Task) bool {
			return tasks[0].Done == false
		}).Return(nil)

		err := app.ToggleDone(&DoneCmd{Id: 1})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("returns error when not found", func(t *testing.T) {
		app, mock := newTestTodo()
		mock.EXPECT().Load(t).Return([]Task{}, nil)

		err := app.ToggleDone(&DoneCmd{Id: 999})
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("returns error on load failure", func(t *testing.T) {
		app, mock := newTestTodo()
		mock.EXPECT().Load(t).Return(nil, fmt.Errorf("error"))

		err := app.ToggleDone(&DoneCmd{Id: 1})
		if err == nil {
			t.Fatal("expected error")
		}
	})
}

// --- DeleteTask ---

func TestDeleteTask(t *testing.T) {
	t.Run("removes task by id", func(t *testing.T) {
		app, mock := newTestTodo()
		mock.EXPECT().Load(t).Return([]Task{{Id: 1, Title: "Task 1"}, {Id: 2, Title: "Task 2"}}, nil)
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

	t.Run("returns error on load failure", func(t *testing.T) {
		app, mock := newTestTodo()
		mock.EXPECT().Load(t).Return(nil, fmt.Errorf("error"))

		err := app.DeleteTask(&DeleteCmd{Id: 1})
		if err == nil {
			t.Fatal("expected error")
		}
	})
}

// --- UpdateTasks ---

func TestUpdateTasks(t *testing.T) {
	t.Run("updates title only", func(t *testing.T) {
		app, mock := newTestTodo()
		mock.EXPECT().Load(t).Return([]Task{{Id: 1, Title: "Old", Description: "desc", Priority: "low"}}, nil)
		mock.EXPECT().Save(t).MatchTasks(func(tasks []Task) bool {
			return tasks[0].Title == "New" && tasks[0].Description == "desc" && tasks[0].Priority == "low"
		}).Return(nil)

		err := app.UpdateTasks(&UpdateCmd{Id: 1, Title: "New"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("updates all fields", func(t *testing.T) {
		app, mock := newTestTodo()
		mock.EXPECT().Load(t).Return([]Task{{Id: 1, Title: "Old", Description: "old", Priority: "low"}}, nil)
		mock.EXPECT().Save(t).MatchTasks(func(tasks []Task) bool {
			return tasks[0].Title == "New" && tasks[0].Description == "new desc" && tasks[0].Priority == "high"
		}).Return(nil)

		err := app.UpdateTasks(&UpdateCmd{Id: 1, Title: "New", Description: "new desc", Priority: "high"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("returns error when not found", func(t *testing.T) {
		app, mock := newTestTodo()
		mock.EXPECT().Load(t).Return([]Task{}, nil)

		err := app.UpdateTasks(&UpdateCmd{Id: 999, Title: "title"})
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("returns error on load failure", func(t *testing.T) {
		app, mock := newTestTodo()
		mock.EXPECT().Load(t).Return(nil, fmt.Errorf("error"))

		err := app.UpdateTasks(&UpdateCmd{Id: 1, Title: "title"})
		if err == nil {
			t.Fatal("expected error")
		}
	})
}

// --- ListTasks ---

func TestListTasks(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		app, mock := newTestTodo()
		mock.EXPECT().Load(t).Return([]Task{}, nil)

		err := app.ListTasks(&ListCmd{})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("loads tasks without error", func(t *testing.T) {
		app, mock := newTestTodo()
		mock.EXPECT().Load(t).Return([]Task{
			{Id: 1, Title: "Task 1", Priority: "low"},
			{Id: 2, Title: "Task 2", Priority: "high", Done: true},
		}, nil)

		err := app.ListTasks(&ListCmd{})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("returns error on load failure", func(t *testing.T) {
		app, mock := newTestTodo()
		mock.EXPECT().Load(t).Return(nil, fmt.Errorf("error"))

		err := app.ListTasks(&ListCmd{})
		if err == nil {
			t.Fatal("expected error")
		}
	})
}
