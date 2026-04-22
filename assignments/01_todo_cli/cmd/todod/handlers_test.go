package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	todo "github.com/trungnc90/learn-golang/assignments/01_todo_cli"
)

// newTestMux creates a mux backed by an in-memory store.
func newTestMux() *http.ServeMux {
	manager := todo.New() // defaults to MemoryStore
	return newMux(manager)
}

// --- POST /tasks ---

func TestHandleAddTask(t *testing.T) {
	t.Run("creates a task", func(t *testing.T) {
		mux := newTestMux()
		req := httptest.NewRequest("POST", "/tasks", strings.NewReader(`{"title":"Buy groceries","description":"Milk","priority":"high"}`))
		w := httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
		}
		var task todo.Task
		json.NewDecoder(w.Body).Decode(&task)
		if task.Title != "Buy groceries" {
			t.Fatalf("expected 'Buy groceries', got '%s'", task.Title)
		}
		if task.Priority != "high" {
			t.Fatalf("expected 'high', got '%s'", task.Priority)
		}
		if task.Id != 1 {
			t.Fatalf("expected id 1, got %d", task.Id)
		}
	})

	t.Run("defaults priority to low", func(t *testing.T) {
		mux := newTestMux()
		req := httptest.NewRequest("POST", "/tasks", strings.NewReader(`{"title":"Task"}`))
		w := httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		var task todo.Task
		json.NewDecoder(w.Body).Decode(&task)
		if task.Priority != "low" {
			t.Fatalf("expected 'low', got '%s'", task.Priority)
		}
	})

	t.Run("returns 400 when title is empty", func(t *testing.T) {
		mux := newTestMux()
		req := httptest.NewRequest("POST", "/tasks", strings.NewReader(`{"description":"no title"}`))
		w := httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", w.Code)
		}
	})

	t.Run("returns 400 for invalid JSON", func(t *testing.T) {
		mux := newTestMux()
		req := httptest.NewRequest("POST", "/tasks", strings.NewReader("not json"))
		w := httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", w.Code)
		}
	})
}

// --- GET /tasks ---

func TestHandleListTasks(t *testing.T) {
	t.Run("returns empty list", func(t *testing.T) {
		mux := newTestMux()
		req := httptest.NewRequest("GET", "/tasks", nil)
		w := httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
		var tasks []todo.Task
		json.NewDecoder(w.Body).Decode(&tasks)
		if len(tasks) != 0 {
			t.Fatalf("expected 0 tasks, got %d", len(tasks))
		}
	})

	t.Run("returns added tasks", func(t *testing.T) {
		mux := newTestMux()
		// Add a task
		mux.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("POST", "/tasks", strings.NewReader(`{"title":"Task 1"}`)))

		req := httptest.NewRequest("GET", "/tasks", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		var tasks []todo.Task
		json.NewDecoder(w.Body).Decode(&tasks)
		if len(tasks) != 1 {
			t.Fatalf("expected 1 task, got %d", len(tasks))
		}
	})

	t.Run("filters done tasks", func(t *testing.T) {
		mux := newTestMux()
		mux.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("POST", "/tasks", strings.NewReader(`{"title":"Pending"}`)))
		mux.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("POST", "/tasks", strings.NewReader(`{"title":"Done"}`)))
		mux.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("PATCH", "/tasks/2/toggle", nil))

		req := httptest.NewRequest("GET", "/tasks?filter=done", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		var tasks []todo.Task
		json.NewDecoder(w.Body).Decode(&tasks)
		if len(tasks) != 1 || tasks[0].Title != "Done" {
			t.Fatalf("expected 1 done task")
		}
	})

	t.Run("filters pending tasks", func(t *testing.T) {
		mux := newTestMux()
		mux.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("POST", "/tasks", strings.NewReader(`{"title":"Pending"}`)))
		mux.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("POST", "/tasks", strings.NewReader(`{"title":"Done"}`)))
		mux.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("PATCH", "/tasks/2/toggle", nil))

		req := httptest.NewRequest("GET", "/tasks?filter=pending", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		var tasks []todo.Task
		json.NewDecoder(w.Body).Decode(&tasks)
		if len(tasks) != 1 || tasks[0].Title != "Pending" {
			t.Fatalf("expected 1 pending task")
		}
	})
}

// --- PATCH /tasks/{id}/toggle ---

func TestHandleToggleDone(t *testing.T) {
	t.Run("toggles task status", func(t *testing.T) {
		mux := newTestMux()
		mux.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("POST", "/tasks", strings.NewReader(`{"title":"Test"}`)))

		req := httptest.NewRequest("PATCH", "/tasks/1/toggle", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
		var task todo.Task
		json.NewDecoder(w.Body).Decode(&task)
		if !task.Done {
			t.Fatal("expected done=true")
		}
	})

	t.Run("toggles back to pending", func(t *testing.T) {
		mux := newTestMux()
		mux.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("POST", "/tasks", strings.NewReader(`{"title":"Test"}`)))
		mux.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("PATCH", "/tasks/1/toggle", nil))

		req := httptest.NewRequest("PATCH", "/tasks/1/toggle", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		var task todo.Task
		json.NewDecoder(w.Body).Decode(&task)
		if task.Done {
			t.Fatal("expected done=false")
		}
	})

	t.Run("returns 404 for missing task", func(t *testing.T) {
		mux := newTestMux()
		req := httptest.NewRequest("PATCH", "/tasks/999/toggle", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", w.Code)
		}
	})
}

// --- DELETE /tasks/{id} ---

func TestHandleDeleteTask(t *testing.T) {
	t.Run("deletes a task", func(t *testing.T) {
		mux := newTestMux()
		mux.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("POST", "/tasks", strings.NewReader(`{"title":"Delete me"}`)))

		req := httptest.NewRequest("DELETE", "/tasks/1", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}

		// Verify it's gone
		listW := httptest.NewRecorder()
		mux.ServeHTTP(listW, httptest.NewRequest("GET", "/tasks", nil))
		var tasks []todo.Task
		json.NewDecoder(listW.Body).Decode(&tasks)
		if len(tasks) != 0 {
			t.Fatalf("expected 0 tasks, got %d", len(tasks))
		}
	})

	t.Run("returns 404 for missing task", func(t *testing.T) {
		mux := newTestMux()
		req := httptest.NewRequest("DELETE", "/tasks/999", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", w.Code)
		}
	})
}

// --- PUT /tasks/{id} ---

func TestHandleUpdateTask(t *testing.T) {
	t.Run("updates a task", func(t *testing.T) {
		mux := newTestMux()
		mux.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("POST", "/tasks", strings.NewReader(`{"title":"Old","priority":"low"}`)))

		req := httptest.NewRequest("PUT", "/tasks/1", strings.NewReader(`{"title":"New","priority":"high"}`))
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
		var task todo.Task
		json.NewDecoder(w.Body).Decode(&task)
		if task.Title != "New" {
			t.Fatalf("expected 'New', got '%s'", task.Title)
		}
		if task.Priority != "high" {
			t.Fatalf("expected 'high', got '%s'", task.Priority)
		}
	})

	t.Run("returns 404 for missing task", func(t *testing.T) {
		mux := newTestMux()
		req := httptest.NewRequest("PUT", "/tasks/999", strings.NewReader(`{"title":"Nope"}`))
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", w.Code)
		}
	})

	t.Run("returns 400 for invalid JSON", func(t *testing.T) {
		mux := newTestMux()
		mux.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("POST", "/tasks", strings.NewReader(`{"title":"Task"}`)))

		req := httptest.NewRequest("PUT", "/tasks/1", strings.NewReader("bad json"))
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", w.Code)
		}
	})
}
