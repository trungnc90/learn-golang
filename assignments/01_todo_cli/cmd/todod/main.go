package main

import (
	"fmt"
	"log"
	"net/http"

	todo "github.com/trungnc90/learn-golang/assignments/01_todo_cli"
	"github.com/trungnc90/learn-golang/assignments/01_todo_cli/infra"
	"github.com/trungnc90/learn-golang/assignments/01_todo_cli/internal/handler"
)

func main() {
	fs := infra.NewFileStore("tasks.json")
	manager := todo.New(todo.WithStorer(fs))
	mux := handler.NewMux(manager)

	fmt.Println("Todo API running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
