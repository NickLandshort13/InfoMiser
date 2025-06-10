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

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
