// How to run tests:
//   go test -v ./...                    Run all tests
//   go test -v -run TestAddTask ./...   Run addTask subtests only
//   go test -count=1 ./...              Run without cache

package todo

import (
	"fmt"
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

// --- AddTask ---

func TestAddTask(t *testing.T) {
	t.Run("calls Load then Save with new task", func(t *testing.T) {
		mock := testStorer()
		mock.EXPECT().Load(t).Return([]Task{}, nil)
		mock.EXPECT().Save(t).MatchTasks(func(tasks []Task) bool {
			return len(tasks) == 1 &&
				tasks[0].Title == "Buy groceries" &&
				tasks[0].Description == "Milk and eggs" &&
				tasks[0].Priority == "high" &&
				tasks[0].Id == 1 &&
				!tasks[0].Done
		}).Return(nil)

		AddTask(mock, &AddCmd{Title: "Buy groceries", Description: "Milk and eggs", Priority: "high"})
	})

	t.Run("default priority is low", func(t *testing.T) {
		mock := testStorer()
		mock.EXPECT().Load(t).Return([]Task{}, nil)
		mock.EXPECT().Save(t).MatchTasks(func(tasks []Task) bool {
			return tasks[0].Priority == "low"
		}).Return(nil)

		AddTask(mock, &AddCmd{Title: "Task"})
	})

	t.Run("appends to existing tasks", func(t *testing.T) {
		existing := []Task{{Id: 3, Title: "Existing"}}
		mock := testStorer()
		mock.EXPECT().Load(t).Return(existing, nil)
		mock.EXPECT().Save(t).MatchTasks(func(tasks []Task) bool {
			return len(tasks) == 2 && tasks[1].Id == 4
		}).Return(nil)

		AddTask(mock, &AddCmd{Title: "New", Priority: "high"})
	})

	t.Run("does not save on load error", func(t *testing.T) {
		mock := testStorer()
		mock.EXPECT().Load(t).Return(nil, fmt.Errorf("disk error"))
		// No Save EXPECT — if AddTask calls Save, it panics

		AddTask(mock, &AddCmd{Title: "Test"})
	})

	t.Run("handles save error", func(t *testing.T) {
		mock := testStorer()
		mock.EXPECT().Load(t).Return([]Task{}, nil)
		mock.EXPECT().Save(t).Return(fmt.Errorf("write failed"))

		AddTask(mock, &AddCmd{Title: "Test"})
	})
}

// --- ToggleDone ---

func TestToggleDone(t *testing.T) {
	t.Run("marks task as done", func(t *testing.T) {
		tasks := []Task{{Id: 1, Title: "Task", Done: false}}
		mock := testStorer()
		mock.EXPECT().Load(t).Return(tasks, nil)
		mock.EXPECT().Save(t).MatchTasks(func(tasks []Task) bool {
			return tasks[0].Done == true
		}).Return(nil)

		ToggleDone(mock, &DoneCmd{Id: 1})
	})

	t.Run("marks task as pending", func(t *testing.T) {
		tasks := []Task{{Id: 1, Title: "Task", Done: true}}
		mock := testStorer()
		mock.EXPECT().Load(t).Return(tasks, nil)
		mock.EXPECT().Save(t).MatchTasks(func(tasks []Task) bool {
			return tasks[0].Done == false
		}).Return(nil)

		ToggleDone(mock, &DoneCmd{Id: 1})
	})

	t.Run("not found does not save", func(t *testing.T) {
		mock := testStorer()
		mock.EXPECT().Load(t).Return([]Task{}, nil)

		ToggleDone(mock, &DoneCmd{Id: 999})
	})

	t.Run("does not save on load error", func(t *testing.T) {
		mock := testStorer()
		mock.EXPECT().Load(t).Return(nil, fmt.Errorf("error"))

		ToggleDone(mock, &DoneCmd{Id: 1})
	})
}

// --- DeleteTask ---

func TestDeleteTask(t *testing.T) {
	t.Run("removes task by id", func(t *testing.T) {
		tasks := []Task{
			{Id: 1, Title: "Task 1"},
			{Id: 2, Title: "Task 2"},
		}
		mock := testStorer()
		mock.EXPECT().Load(t).Return(tasks, nil)
		mock.EXPECT().Save(t).MatchTasks(func(tasks []Task) bool {
			return len(tasks) == 1 && tasks[0].Id == 2
		}).Return(nil)

		DeleteTask(mock, &DeleteCmd{Id: 1})
	})

	t.Run("not found does not save", func(t *testing.T) {
		mock := testStorer()
		mock.EXPECT().Load(t).Return([]Task{}, nil)

		DeleteTask(mock, &DeleteCmd{Id: 999})
	})

	t.Run("does not save on load error", func(t *testing.T) {
		mock := testStorer()
		mock.EXPECT().Load(t).Return(nil, fmt.Errorf("error"))

		DeleteTask(mock, &DeleteCmd{Id: 1})
	})
}

// --- UpdateTasks ---

func TestUpdateTasks(t *testing.T) {
	t.Run("updates title only", func(t *testing.T) {
		tasks := []Task{{Id: 1, Title: "Old", Description: "desc", Priority: "low"}}
		mock := testStorer()
		mock.EXPECT().Load(t).Return(tasks, nil)
		mock.EXPECT().Save(t).MatchTasks(func(tasks []Task) bool {
			return tasks[0].Title == "New" &&
				tasks[0].Description == "desc" &&
				tasks[0].Priority == "low"
		}).Return(nil)

		UpdateTasks(mock, &UpdateCmd{Id: 1, Title: "New"})
	})

	t.Run("updates all fields", func(t *testing.T) {
		tasks := []Task{{Id: 1, Title: "Old", Description: "old", Priority: "low"}}
		mock := testStorer()
		mock.EXPECT().Load(t).Return(tasks, nil)
		mock.EXPECT().Save(t).MatchTasks(func(tasks []Task) bool {
			return tasks[0].Title == "New" &&
				tasks[0].Description == "new desc" &&
				tasks[0].Priority == "high"
		}).Return(nil)

		UpdateTasks(mock, &UpdateCmd{Id: 1, Title: "New", Description: "new desc", Priority: "high"})
	})

	t.Run("not found does not save", func(t *testing.T) {
		mock := testStorer()
		mock.EXPECT().Load(t).Return([]Task{}, nil)

		UpdateTasks(mock, &UpdateCmd{Id: 999, Title: "title"})
	})

	t.Run("does not save on load error", func(t *testing.T) {
		mock := testStorer()
		mock.EXPECT().Load(t).Return(nil, fmt.Errorf("error"))

		UpdateTasks(mock, &UpdateCmd{Id: 1, Title: "title"})
	})
}

// --- ListTasks ---

func TestListTasks(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		mock := testStorer()
		mock.EXPECT().Load(t).Return([]Task{}, nil)

		ListTasks(mock, &ListCmd{})
	})

	t.Run("loads tasks without error", func(t *testing.T) {
		tasks := []Task{
			{Id: 1, Title: "Task 1", Priority: "low"},
			{Id: 2, Title: "Task 2", Priority: "high", Done: true},
		}
		mock := testStorer()
		mock.EXPECT().Load(t).Return(tasks, nil)

		ListTasks(mock, &ListCmd{})
	})

	t.Run("handles load error", func(t *testing.T) {
		mock := testStorer()
		mock.EXPECT().Load(t).Return(nil, fmt.Errorf("error"))

		ListTasks(mock, &ListCmd{})
	})
}
