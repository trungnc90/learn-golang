package main

import (
	"log"

	todo "github.com/trungnc90/learn-golang/assignments/01_todo"
	"github.com/trungnc90/learn-golang/assignments/01_todo/infra"
	"github.com/trungnc90/learn-golang/assignments/01_todo/rest"
)

func main() {
	// storer := infra.NewFileStore("tasks.json")

	storer, err := infra.NewPostgresStore("postgres://todouser:todopass@localhost:5432/tododb?sslmode=disable")
	if err != nil {
		log.Fatalf("connect to database: %v", err)
	}
	defer storer.Close()

	manager := todo.NewManager(storer)
	server := rest.NewServer(manager)

	log.Fatal(server.Run(":8080"))
}
