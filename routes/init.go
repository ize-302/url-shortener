// Package routes
package routes

import (
	"log"
	"net/http"
)

func Requests() {
	http.HandleFunc("/urls", urlHandler)
	http.HandleFunc("/{id}", gotoHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
