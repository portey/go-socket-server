package main

import (
	"./server"
	"log"
	"net/http"
)

func main() {
	serverInstance := server.NewServer("/")
	go serverInstance.Listen()

	log.Fatal(http.ListenAndServe(":8086", nil))
}
