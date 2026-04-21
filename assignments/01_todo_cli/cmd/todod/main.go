package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	todo "github.com/trungnc90/learn-golang/assignments/01_todo_cli"
	"github.com/trungnc90/learn-golang/assignments/01_todo_cli/infra"
)

var manager *todo.Todo

func main() {
	fs := infra.NewFileStore("tasks.json")
	manager = todo.New(todo.WithStorer(fs))

	http.HandleFunc("/tasks", handleTasks)
	http.HandleFunc("/tasks/", handleTaskByID)

	fmt.Println("Todo API running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// handleTasks handles GET /tasks and POST /tasks
func handleTasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		filter := r.URL.Query().Get("filter")
		tasks, err := manager.ListTasks(&todo.ListCmd{Filter: filter})
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, tasks)

	case http.MethodPost:
		var cmd todo.AddCmd
		if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}
		if cmd.Title == "" {
			writeError(w, http.StatusBadRequest, "title is required")
			return
		}
		task, err := manager.AddTask(&cmd)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusCreated, task)

	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// handleTaskByID handles GET/PUT/DELETE /tasks/{id} and PATCH /tasks/{id}/toggle
func handleTaskByID(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/tasks/")
	parts := strings.Split(path, "/")

	id, err := strconv.Atoi(parts[0])
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid task id")
		return
	}

	// PATCH /tasks/{id}/toggle
	if len(parts) == 2 && parts[1] == "toggle" && r.Method == http.MethodPatch {
		task, err := manager.ToggleDone(&todo.DoneCmd{Id: id})
		if err != nil {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, task)
		return
	}

	switch r.Method {
	case http.MethodPut:
		var cmd todo.UpdateCmd
		if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}
		cmd.Id = id
		task, err := manager.UpdateTasks(&cmd)
		if err != nil {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, task)

	case http.MethodDelete:
		if err := manager.DeleteTask(&todo.DeleteCmd{Id: id}); err != nil {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"message": fmt.Sprintf("task #%d deleted", id)})

	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
