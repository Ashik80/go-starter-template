package form

import (
	"html/template"
	"net/http"

	"github.com/gorilla/csrf"
)

type Form struct {
	CSRF template.HTML
}

func NewForm(r *http.Request) Form {
	return Form{
		CSRF: csrf.TemplateField(r),
	}
}
