package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	todo "github.com/trungnc90/learn-golang/assignments/01_todo_cli"
	"github.com/trungnc90/learn-golang/assignments/01_todo_cli/infra"
)

var manager *todo.Todo

func main() {
	fs := infra.NewFileStore("tasks.json")
	manager = todo.New(todo.WithStorer(fs))

	http.HandleFunc("GET /tasks", handleListTasks)
	http.HandleFunc("POST /tasks", handleAddTask)
	http.HandleFunc("PUT /tasks/{id}", handleUpdateTask)
	http.HandleFunc("DELETE /tasks/{id}", handleDeleteTask)
	http.HandleFunc("PATCH /tasks/{id}/toggle", handleToggleDone)

	fmt.Println("Todo API running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleListTasks(w http.ResponseWriter, r *http.Request) {
	filter := r.URL.Query().Get("filter")
	tasks, err := manager.ListTasks(&todo.ListCmd{Filter: filter})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, tasks)
}

func handleAddTask(w http.ResponseWriter, r *http.Request) {
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
}

func handleUpdateTask(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid task id")
		return
	}
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
}

func handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid task id")
		return
	}
	if err := manager.DeleteTask(&todo.DeleteCmd{Id: id}); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": fmt.Sprintf("task #%d deleted", id)})
}

func handleToggleDone(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid task id")
		return
	}
	task, err := manager.ToggleDone(&todo.DoneCmd{Id: id})
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, task)
}

func parseID(r *http.Request) (int, error) {
	idStr := r.PathValue("id")
	var id int
	_, err := fmt.Sscanf(idStr, "%d", &id)
	return id, err
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
