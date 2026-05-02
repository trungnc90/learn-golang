// How to run tests:
//   go test -v ./...                    Run all tests
//   go test -v -run TestAddTask ./...   Run addTask subtests only
//   go test -count=1 ./...              Run without cache

package todo

import (
	"fmt"
	"testing"
	"time"
)

// newTestTodo creates a Todo with a mock storer for testing.
func newTestTodo() (*Todo, *storer) {
	mock := testStorer()
	return NewManager(mock), mock
}

// --- AddTask ---

func TestAddTask(t *testing.T) {
	t.Run("returns created task", func(t *testing.T) {
		app, mock := newTestTodo()
		mock.EXPECT().Create(t).MatchTask(func(task Task) bool {
			return task.Title == "Buy groceries" && task.Priority == "high"
		}).Return(Task{Id: 1, Title: "Buy groceries", Description: "Milk and eggs", Priority: "high", CreatedAt: time.Now()}, nil)

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
		mock.EXPECT().Create(t).MatchTask(func(task Task) bool {
			return task.Priority == "low"
		}).Return(Task{Id: 1, Title: "Task", Priority: "low", CreatedAt: time.Now()}, nil)

		task, err := app.AddTask(&AddCmd{Title: "Task"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if task.Priority != "low" {
			t.Fatalf("got '%s'", task.Priority)
		}
	})

	t.Run("returns error on create failure", func(t *testing.T) {
		app, mock := newTestTodo()
		mock.EXPECT().Create(t).Return(Task{}, fmt.Errorf("db error"))

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
		mock.EXPECT().List(t).WithFilter("").Return([]Task{
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

	t.Run("passes done filter", func(t *testing.T) {
		app, mock := newTestTodo()
		mock.EXPECT().List(t).WithFilter("done").Return([]Task{
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

	t.Run("passes pending filter", func(t *testing.T) {
		app, mock := newTestTodo()
		mock.EXPECT().List(t).WithFilter("pending").Return([]Task{
			{Id: 1, Title: "Pending", Done: false},
		}, nil)

		tasks, err := app.ListTasks(&ListCmd{Filter: "pending"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(tasks) != 1 || tasks[0].Title != "Pending" {
			t.Fatalf("expected only pending task")
		}
	})

	t.Run("returns error on list failure", func(t *testing.T) {
		app, mock := newTestTodo()
		mock.EXPECT().List(t).Return(nil, fmt.Errorf("error"))

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
		mock.EXPECT().ToggleDone(t).WithID(1).Return(Task{Id: 1, Done: true}, nil)

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
		mock.EXPECT().ToggleDone(t).Return(Task{}, fmt.Errorf("task #999 not found"))

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
		mock.EXPECT().Delete(t).WithID(1).Return(nil)

		err := app.DeleteTask(&DeleteCmd{Id: 1})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("returns error when not found", func(t *testing.T) {
		app, mock := newTestTodo()
		mock.EXPECT().Delete(t).Return(fmt.Errorf("task #999 not found"))

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
		mock.EXPECT().Update(t).MatchTask(func(task Task) bool {
			return task.Id == 1 && task.Title == "New" && task.Priority == "high"
		}).Return(Task{Id: 1, Title: "New", Description: "desc", Priority: "high"}, nil)

		task, err := app.UpdateTasks(&UpdateCmd{Id: 1, Title: "New", Priority: "high"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if task.Title != "New" {
			t.Fatalf("got '%s'", task.Title)
		}
		if task.Priority != "high" {
			t.Fatalf("got '%s'", task.Priority)
		}
	})

	t.Run("returns error when not found", func(t *testing.T) {
		app, mock := newTestTodo()
		mock.EXPECT().Update(t).Return(Task{}, fmt.Errorf("task #999 not found"))

		_, err := app.UpdateTasks(&UpdateCmd{Id: 999, Title: "title"})
		if err == nil {
			t.Fatal("expected error")
		}
	})
}
