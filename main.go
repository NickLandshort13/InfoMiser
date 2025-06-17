package main

import (
	"html/template"
	"net/http"

	"infomiser/internal/handlers"
)

func main() {
	templates := template.Must(template.ParseGlob("templates/*.html"))
	h := handlers.New(templates)

	http.HandleFunc("/", h.Home)
	http.HandleFunc("/lookup", h.Lookup)

	println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
