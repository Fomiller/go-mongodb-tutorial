package main

import (
	"net/http"

	"github.com/fomiller/go-mongodb-tutorial/API"
)

func main() {
	// Route Handlers
	http.HandleFunc("/", API.IndexHandler)
	http.HandleFunc("/api/create", API.CreateHandler)
	http.HandleFunc("/api/createmany", API.CreateManyHandler)
	http.HandleFunc("/api/update", API.UpdateHandler)
	http.HandleFunc("/api/find", API.FindHandler)
	http.HandleFunc("/api/findmany", API.FindManyHandler)
	http.HandleFunc("/api/delete", API.DeleteHandler)
	// handle favicon
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":8080", nil)
}
