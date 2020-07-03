package main

import (
	"net/http"

	"github.com/fomiller/go-mongodb-tutorial/API"
)

func main() {
	http.HandleFunc("/", API.IndexHandler)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":8080", nil)
}
