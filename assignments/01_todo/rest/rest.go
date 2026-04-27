package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	todo "github.com/trungnc90/learn-golang/assignments/01_todo"
)

//go:generate go-mock-gen --struct=Server
type Server struct {
	manager todo.Manager
}

func NewServer(manager todo.Manager) *Server {
	return &Server{manager: manager}
}

type route struct {
	pattern string
	handler http.HandlerFunc
}

func (s *Server) Run(addr string) error {
	mux := http.NewServeMux()

	routes := []route{
		{"GET /tasks", s.HandleListTasks},
		{"POST /tasks", s.HandleAddTask},
		{"PUT /tasks/{id}", s.HandleUpdateTask},
		{"DELETE /tasks/{id}", s.HandleDeleteTask},
		{"PATCH /tasks/{id}/toggle", s.HandleToggleDone},
	}

	for _, r := range routes {
		mux.HandleFunc(r.pattern, r.handler)
	}

	fmt.Printf("Todo API running on http://%s\n", addr)
	return http.ListenAndServe(addr, mux)
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
