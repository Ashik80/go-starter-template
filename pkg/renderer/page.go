package renderer

import "html/template"

var pages map[string]*template.Template

func RegisterPageTemplates() {
	pages = map[string]*template.Template{
		"home":         ParseTemplate("home"),
		"todos":        ParseTemplate("todos"),
		"todo-details": ParseTemplate("todo-details"),
		"login":        ParseTemplate("login", "auth"),
		"signup":       ParseTemplate("signup", "auth"),
	}
}

func GetPageTemplate(name string) *template.Template {
	return pages[name]
}
