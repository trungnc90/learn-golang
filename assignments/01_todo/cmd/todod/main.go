package main

import (
	"log"

	todo "github.com/trungnc90/learn-golang/assignments/01_todo"
	"github.com/trungnc90/learn-golang/assignments/01_todo/infra"
	"github.com/trungnc90/learn-golang/assignments/01_todo/rest"
)

func main() {
	storer := infra.NewFileStore("tasks.json")
	manager := todo.NewManager(storer)
	server := rest.NewServer(manager)

	log.Fatal(server.Run(":8080"))
}
