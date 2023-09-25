package main

import (
	"net/http"
	"simple-http-server/api"
)

func main() {
	server := api.NewServer()
	http.ListenAndServe(":8080", server)
}
