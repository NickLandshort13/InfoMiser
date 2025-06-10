package handlers

import "html/template"

type Handlers struct {
    templates *template.Template
}

func New(tmpl *template.Template) *Handlers {
    return &Handlers{templates: tmpl}
}