package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	todo "github.com/trungnc90/learn-golang/assignments/01_todo_cli"
)

//go:generate go-mock-gen -i TaskManager
type TaskManager interface {
	AddTask(cmd *todo.AddCmd) (*todo.Task, error)
	ListTasks(cmd *todo.ListCmd) ([]todo.Task, error)
	UpdateTasks(cmd *todo.UpdateCmd) (*todo.Task, error)
	DeleteTask(cmd *todo.DeleteCmd) error
	ToggleDone(cmd *todo.DoneCmd) (*todo.Task, error)
}

type Server struct {
	manager TaskManager
}

func NewMux(manager TaskManager) *http.ServeMux {
	s := &Server{manager: manager}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /tasks", s.HandleListTasks)
	mux.HandleFunc("POST /tasks", s.HandleAddTask)
	mux.HandleFunc("PUT /tasks/{id}", s.HandleUpdateTask)
	mux.HandleFunc("DELETE /tasks/{id}", s.HandleDeleteTask)
	mux.HandleFunc("PATCH /tasks/{id}/toggle", s.HandleToggleDone)
	return mux
}

func (s *Server) HandleListTasks(w http.ResponseWriter, r *http.Request) {
	filter := r.URL.Query().Get("filter")
	tasks, err := s.manager.ListTasks(&todo.ListCmd{Filter: filter})
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, tasks)
}

func (s *Server) HandleAddTask(w http.ResponseWriter, r *http.Request) {
	var cmd todo.AddCmd
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if cmd.Title == "" {
		WriteError(w, http.StatusBadRequest, "title is required")
		return
	}
	task, err := s.manager.AddTask(&cmd)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusCreated, task)
}

func (s *Server) HandleUpdateTask(w http.ResponseWriter, r *http.Request) {
	id, err := ParseID(r)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid task id")
		return
	}
	var cmd todo.UpdateCmd
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	cmd.Id = id
	task, err := s.manager.UpdateTasks(&cmd)
	if err != nil {
		WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, task)
}

func (s *Server) HandleDeleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := ParseID(r)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid task id")
		return
	}
	if err := s.manager.DeleteTask(&todo.DeleteCmd{Id: id}); err != nil {
		WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, map[string]string{"message": fmt.Sprintf("task #%d deleted", id)})
}

func (s *Server) HandleToggleDone(w http.ResponseWriter, r *http.Request) {
	id, err := ParseID(r)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid task id")
		return
	}
	task, err := s.manager.ToggleDone(&todo.DoneCmd{Id: id})
	if err != nil {
		WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, task)
}

func ParseID(r *http.Request) (int, error) {
	idStr := r.PathValue("id")
	var id int
	_, err := fmt.Sscanf(idStr, "%d", &id)
	return id, err
}

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, map[string]string{"error": message})
}
