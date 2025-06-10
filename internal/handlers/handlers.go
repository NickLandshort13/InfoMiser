package handlers

import (
	"html/template"
	"net/http"
)

type Handlers struct {
	templates *template.Template
}

func New(tmpl *template.Template) *Handlers {
	return &Handlers{templates: tmpl}
}

func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	h.templates.ExecuteTemplate(w, "index.html", map[string]string{
		"Title": "InfoMiser — OSINT Lookup",
	})
}
