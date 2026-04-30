package main

import (
	"log"

	"github.com/floffah/catena/internal/app/api"
)

func main() {
	server := api.NewServer()

	log.Fatal(server.ListenAndServe())
}
