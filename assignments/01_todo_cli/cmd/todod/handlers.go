package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	todo "github.com/trungnc90/learn-golang/assignments/01_todo_cli"
)

func (s *server) handleListTasks(w http.ResponseWriter, r *http.Request) {
	filter := r.URL.Query().Get("filter")
	tasks, err := s.manager.ListTasks(&todo.ListCmd{Filter: filter})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, tasks)
}

func (s *server) handleAddTask(w http.ResponseWriter, r *http.Request) {
	var cmd todo.AddCmd
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if cmd.Title == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}
	task, err := s.manager.AddTask(&cmd)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, task)
}

func (s *server) handleUpdateTask(w http.ResponseWriter, r *http.Request) {
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
	task, err := s.manager.UpdateTasks(&cmd)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, task)
}

func (s *server) handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid task id")
		return
	}
	if err := s.manager.DeleteTask(&todo.DeleteCmd{Id: id}); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": fmt.Sprintf("task #%d deleted", id)})
}

func (s *server) handleToggleDone(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid task id")
		return
	}
	task, err := s.manager.ToggleDone(&todo.DoneCmd{Id: id})
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
