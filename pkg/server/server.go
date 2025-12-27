package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

const webDir = "./web"
const defaultPort = "7540"

func Start() {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = defaultPort
	}

	fs := http.FileServer(http.Dir(webDir))
	http.Handle("/", fs)

	address := fmt.Sprintf(":%s", port)
	fmt.Println("Server started on", address)

	log.Fatal(http.ListenAndServe(address, nil))
}
