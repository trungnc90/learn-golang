package main

import (
	"fmt"
	"log"

	"example.com/greetings"
)

func main() {

	// set log prefix and flag 0 (disable printing the time, source file, line number)
	log.SetPrefix("greetings: ")
	log.SetFlags(0)

	message, err := greetings.Hello("Noi")
	// if returing an error, print to console and exit
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(message)
}
