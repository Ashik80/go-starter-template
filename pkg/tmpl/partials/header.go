package partialTmpl

import (
	"html/template"
	"net/http"

	"go-starter-template/pkg/infrastructure/renderer"
)

type Header struct {
	tmpl *template.Template
}

func NewHeader() *Header {
	return &Header{
		tmpl: renderer.GetBaseTemplate(),
	}
}

func (p *Header) Execute(w http.ResponseWriter, data any) error {
	return p.tmpl.ExecuteTemplate(w, "header", data)
}
