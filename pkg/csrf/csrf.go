package csrf

import (
	"html/template"
	"net/http"

	"github.com/gorilla/csrf"
)

func GetCSRFField(r *http.Request) template.HTML {
	return csrf.TemplateField(r)
}
