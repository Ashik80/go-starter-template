package controllers

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/gorilla/csrf"

	"go-starter-template/pkg/application/command"
	"go-starter-template/pkg/application/interfaces"
	"go-starter-template/pkg/application/result"
	"go-starter-template/pkg/infrastructure/config"
	"go-starter-template/pkg/infrastructure/middlewares"
	"go-starter-template/pkg/infrastructure/renderer"
	"go-starter-template/pkg/infrastructure/router"
	"go-starter-template/pkg/utils"
)

type TodoController struct {
	todoService interfaces.TodoService
}

type TodoCreateForm struct {
	CSRF        template.HTML
	ID          int
	Title       string
	Description string
	Error       string
}

type TodoDeleteForm struct {
	CSRF template.HTML
	ID   int
}

func NewTodoController(r router.Router, todoService interfaces.TodoService, sessionService interfaces.SessionService, config *config.Config) {
	controller := &TodoController{
		todoService: todoService,
	}

	authMiddleware := middlewares.AuthMiddleware(config.Env, sessionService)

	r.Route("/todos", func(r router.Router) {
		// INFO: to apply middleware to a group of routes use Use method
		r.Use(authMiddleware)

		r.Get("/", controller.List)
		r.Get("/{id}", controller.Get)
		r.Post("/", controller.Create)
		r.Put("/{id}", controller.Update)
		r.Delete("/{id}", controller.Delete)
	})
}

func (tc *TodoController) List(w http.ResponseWriter, r *http.Request) {
	page := renderer.GetPageTemplate("todos")
	data := map[string]interface{}{
		"Title": "Todos",
		"Path":  "/todos",
		"Form": &TodoCreateForm{
			CSRF:        csrf.TemplateField(r),
			Title:       "",
			Description: "",
			Error:       "",
		},
		"Todos": []*result.TodoResult{},
		"Error": "",
	}
	res, err := tc.todoService.ListTodos(r.Context())
	if err != nil {
		data["Error"] = err.Error()
		w.WriteHeader(500)
		page.ExecuteTemplate(w, "base", data)
		return
	}
	data["Todos"] = res.Todos
	w.WriteHeader(200)
	page.ExecuteTemplate(w, "base", data)
}

func (tc *TodoController) Get(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseToInt(router.GetParam(r, "id"))
	if err != nil {
		w.Header().Add("Hx-Location", "/todos")
		w.WriteHeader(400)
		return
	}

	res, err := tc.todoService.GetTodo(r.Context(), id)
	if err != nil {
		w.Header().Add("Hx-Location", "/todos")
		w.WriteHeader(404)
		return
	}

	page := renderer.GetPageTemplate("todo-details")
	data := map[string]interface{}{
		"Title": res.Todo.Title,
		"Path":  "/todos/" + utils.ParseToString(id),
		"Todo":  res.Todo,
		"EditForm": &TodoCreateForm{
			CSRF:        csrf.TemplateField(r),
			ID:          res.Todo.ID,
			Title:       res.Todo.Title,
			Description: res.Todo.Description,
			Error:       "",
		},
		"DeleteForm": &TodoDeleteForm{
			CSRF: csrf.TemplateField(r),
			ID:   res.Todo.ID,
		},
		"Error": "",
	}
	w.WriteHeader(200)
	page.ExecuteTemplate(w, "base", data)
}

func (tc *TodoController) Create(w http.ResponseWriter, r *http.Request) {
	tmpl := renderer.GetBaseTemplate()
	form := &TodoCreateForm{
		CSRF:        csrf.TemplateField(r),
		Title:       strings.TrimSpace(r.FormValue("title")),
		Description: strings.TrimSpace(r.FormValue("description")),
	}

	result, err := tc.todoService.CreateTodo(r.Context(), &command.CreateTodoCommand{
		Title:       form.Title,
		Description: form.Description,
	})
	if err != nil {
		form.Error = err.Error()
		w.WriteHeader(400)
		tmpl.ExecuteTemplate(w, "todo-create-form", form)
		return
	}

	form = &TodoCreateForm{
		CSRF: csrf.TemplateField(r),
	}

	w.WriteHeader(201)
	tmpl.ExecuteTemplate(w, "todo-item-oob", result.Todo)
	tmpl.ExecuteTemplate(w, "todo-create-form", form)
}

func (tc *TodoController) Update(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseToInt(router.GetParam(r, "id"))
	if err != nil {
		w.Header().Add("Hx-Location", "/todos")
		w.WriteHeader(400)
		return
	}

	tmpl := renderer.GetBaseTemplate()

	editForm := &TodoCreateForm{
		CSRF:        csrf.TemplateField(r),
		ID:          id,
		Title:       strings.TrimSpace(r.FormValue("title")),
		Description: strings.TrimSpace(r.FormValue("description")),
	}

	result, err := tc.todoService.UpdateTodo(r.Context(), &command.UpdateTodoCommand{
		ID:          id,
		Title:       editForm.Title,
		Description: editForm.Description,
	})
	if err != nil {
		editForm.Error = err.Error()
		w.WriteHeader(400)
		tmpl.ExecuteTemplate(w, "todo-edit-form", editForm)
		return
	}

	w.Header().Add("Hx-Trigger", "close_edit_form")
	w.WriteHeader(200)

	editForm = &TodoCreateForm{
		CSRF:        csrf.TemplateField(r),
		ID:          id,
		Title:       result.Todo.Title,
		Description: result.Todo.Description,
	}

	tmpl.ExecuteTemplate(w, "todo-details-info-oob", result.Todo)
	tmpl.ExecuteTemplate(w, "todo-edit-form", editForm)
}

func (tc *TodoController) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseToInt(router.GetParam(r, "id"))
	if err != nil {
		w.Header().Add("Hx-Location", "/todos")
		w.WriteHeader(400)
		return
	}
	deleteForm := &TodoDeleteForm{
		CSRF: csrf.TemplateField(r),
		ID:   id,
	}

	page := renderer.GetPageTemplate("todo-details")
	if err := tc.todoService.DeleteTodo(r.Context(), &command.DeleteTodoCommand{ID: id}); err != nil {
		w.WriteHeader(400)
		page.ExecuteTemplate(w, "todo-delete-form", deleteForm)
		errorString := fmt.Sprintf("<div id=\"error-message\" hx-swap-oob=\"true\"><p style='color: red;'>%s</p></div>", err.Error())
		renderer.RenderString(w, errorString, nil)
		return
	}

	w.Header().Add("Hx-Location", "/todos")
	w.WriteHeader(200)
}
