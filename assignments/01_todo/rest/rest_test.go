package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	todo "github.com/trungnc90/learn-golang/assignments/01_todo"
)

func newTestServer(t *testing.T) (*Server, *manager) {
	mock := testManager()
	return NewServer(mock), mock
}

// --- HandleAddTask ---

func TestHandleAddTask(t *testing.T) {
	t.Run("creates a task", func(t *testing.T) {
		s, mock := newTestServer(t)
		mock.EXPECT().AddTask(t).MatchCmd(func(cmd *todo.AddCmd) bool {
			return cmd.Title == "Buy groceries" && cmd.Priority == "high"
		}).Return(&todo.Task{Id: 1, Title: "Buy groceries", Priority: "high"}, nil)

		req := httptest.NewRequest("POST", "/tasks", strings.NewReader(`{"title":"Buy groceries","priority":"high"}`))
		w := httptest.NewRecorder()
		s.HandleAddTask(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
		}
		var task todo.Task
		json.NewDecoder(w.Body).Decode(&task)
		if task.Title != "Buy groceries" {
			t.Fatalf("expected 'Buy groceries', got '%s'", task.Title)
		}
	})

	t.Run("returns 400 when title is empty", func(t *testing.T) {
		s, _ := newTestServer(t)
		req := httptest.NewRequest("POST", "/tasks", strings.NewReader(`{"description":"no title"}`))
		w := httptest.NewRecorder()
		s.HandleAddTask(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", w.Code)
		}
	})

	t.Run("returns 400 for invalid JSON", func(t *testing.T) {
		s, _ := newTestServer(t)
		req := httptest.NewRequest("POST", "/tasks", strings.NewReader("not json"))
		w := httptest.NewRecorder()
		s.HandleAddTask(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", w.Code)
		}
	})

	t.Run("returns 500 on manager error", func(t *testing.T) {
		s, mock := newTestServer(t)
		mock.EXPECT().AddTask(t).Return(nil, fmt.Errorf("store failed"))

		req := httptest.NewRequest("POST", "/tasks", strings.NewReader(`{"title":"Test"}`))
		w := httptest.NewRecorder()
		s.HandleAddTask(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})
}

// --- HandleListTasks ---

func TestHandleListTasks(t *testing.T) {
	t.Run("returns tasks", func(t *testing.T) {
		s, mock := newTestServer(t)
		mock.EXPECT().ListTasks(t).Return([]todo.Task{
			{Id: 1, Title: "Task 1"},
			{Id: 2, Title: "Task 2"},
		}, nil)

		req := httptest.NewRequest("GET", "/tasks", nil)
		w := httptest.NewRecorder()
		s.HandleListTasks(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
		var tasks []todo.Task
		json.NewDecoder(w.Body).Decode(&tasks)
		if len(tasks) != 2 {
			t.Fatalf("expected 2 tasks, got %d", len(tasks))
		}
	})

	t.Run("passes filter query param", func(t *testing.T) {
		s, mock := newTestServer(t)
		mock.EXPECT().ListTasks(t).MatchCmd(func(cmd *todo.ListCmd) bool {
			return cmd.Filter == "done"
		}).Return([]todo.Task{}, nil)

		req := httptest.NewRequest("GET", "/tasks?filter=done", nil)
		w := httptest.NewRecorder()
		s.HandleListTasks(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
	})

	t.Run("returns 500 on error", func(t *testing.T) {
		s, mock := newTestServer(t)
		mock.EXPECT().ListTasks(t).Return(nil, fmt.Errorf("load failed"))

		req := httptest.NewRequest("GET", "/tasks", nil)
		w := httptest.NewRecorder()
		s.HandleListTasks(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})
}

// --- HandleToggleDone ---

func TestHandleToggleDone(t *testing.T) {
	t.Run("toggles task", func(t *testing.T) {
		s, mock := newTestServer(t)
		mock.EXPECT().ToggleDone(t).MatchCmd(func(cmd *todo.DoneCmd) bool {
			return cmd.Id == 1
		}).Return(&todo.Task{Id: 1, Done: true}, nil)

		req := httptest.NewRequest("PATCH", "/tasks/1/toggle", nil)
		req.SetPathValue("id", "1")
		w := httptest.NewRecorder()
		s.HandleToggleDone(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
		var task todo.Task
		json.NewDecoder(w.Body).Decode(&task)
		if !task.Done {
			t.Fatal("expected done=true")
		}
	})

	t.Run("returns 404 for missing task", func(t *testing.T) {
		s, mock := newTestServer(t)
		mock.EXPECT().ToggleDone(t).Return(nil, fmt.Errorf("not found"))

		req := httptest.NewRequest("PATCH", "/tasks/999/toggle", nil)
		req.SetPathValue("id", "999")
		w := httptest.NewRecorder()
		s.HandleToggleDone(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", w.Code)
		}
	})
}

// --- HandleDeleteTask ---

func TestHandleDeleteTask(t *testing.T) {
	t.Run("deletes a task", func(t *testing.T) {
		s, mock := newTestServer(t)
		mock.EXPECT().DeleteTask(t).MatchCmd(func(cmd *todo.DeleteCmd) bool {
			return cmd.Id == 1
		}).Return(nil)

		req := httptest.NewRequest("DELETE", "/tasks/1", nil)
		req.SetPathValue("id", "1")
		w := httptest.NewRecorder()
		s.HandleDeleteTask(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
	})

	t.Run("returns 404 for missing task", func(t *testing.T) {
		s, mock := newTestServer(t)
		mock.EXPECT().DeleteTask(t).Return(fmt.Errorf("not found"))

		req := httptest.NewRequest("DELETE", "/tasks/999", nil)
		req.SetPathValue("id", "999")
		w := httptest.NewRecorder()
		s.HandleDeleteTask(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", w.Code)
		}
	})
}

// --- HandleUpdateTask ---

func TestHandleUpdateTask(t *testing.T) {
	t.Run("updates a task", func(t *testing.T) {
		s, mock := newTestServer(t)
		mock.EXPECT().UpdateTasks(t).MatchCmd(func(cmd *todo.UpdateCmd) bool {
			return cmd.Id == 1 && cmd.Title == "New"
		}).Return(&todo.Task{Id: 1, Title: "New", Priority: "high"}, nil)

		req := httptest.NewRequest("PUT", "/tasks/1", strings.NewReader(`{"title":"New","priority":"high"}`))
		req.SetPathValue("id", "1")
		w := httptest.NewRecorder()
		s.HandleUpdateTask(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
		var task todo.Task
		json.NewDecoder(w.Body).Decode(&task)
		if task.Title != "New" {
			t.Fatalf("expected 'New', got '%s'", task.Title)
		}
	})

	t.Run("returns 404 for missing task", func(t *testing.T) {
		s, mock := newTestServer(t)
		mock.EXPECT().UpdateTasks(t).Return(nil, fmt.Errorf("not found"))

		req := httptest.NewRequest("PUT", "/tasks/999", strings.NewReader(`{"title":"Nope"}`))
		req.SetPathValue("id", "999")
		w := httptest.NewRecorder()
		s.HandleUpdateTask(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", w.Code)
		}
	})

	t.Run("returns 400 for invalid JSON", func(t *testing.T) {
		s, _ := newTestServer(t)
		req := httptest.NewRequest("PUT", "/tasks/1", strings.NewReader("bad json"))
		req.SetPathValue("id", "1")
		w := httptest.NewRecorder()
		s.HandleUpdateTask(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", w.Code)
		}
	})
}
