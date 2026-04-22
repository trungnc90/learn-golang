package main

import (
	"net/http"

	todo "github.com/trungnc90/learn-golang/assignments/01_todo_cli"
)

type server struct {
	manager *todo.Todo
}

func newMux(manager *todo.Todo) *http.ServeMux {
	s := &server{manager: manager}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /tasks", s.handleListTasks)
	mux.HandleFunc("POST /tasks", s.handleAddTask)
	mux.HandleFunc("PUT /tasks/{id}", s.handleUpdateTask)
	mux.HandleFunc("DELETE /tasks/{id}", s.handleDeleteTask)
	mux.HandleFunc("PATCH /tasks/{id}/toggle", s.handleToggleDone)
	return mux
}
