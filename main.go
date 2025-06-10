package main

import (
	"html/template"
	"infomiser/internal/handlers"
	"net/http"
)

func main() {
	templates := template.Must(template.ParseGlob("templates/*.html"))
	h := handlers.New(templates)

	http.HandleFunc("/", h.Home)
	http.HandleFunc("/lookup", h.Lookup)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
