package main

import (
	"log"

	"github.com/floffah/catena/internal/app/api"
)

func main() {
	server, err := api.NewServer()
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(server.ListenAndServe())
}
