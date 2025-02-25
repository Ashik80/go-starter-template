package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"go-starter-template/pkg/app"
	"go-starter-template/pkg/infrastructure"
	"go-starter-template/pkg/middlewares"
	"go-starter-template/pkg/page"
	"go-starter-template/pkg/service"
)

type TodoHandler struct {
	infrastructure.Router
	infrastructure.TemplateRenderer
	todoService    service.TodoService
	authMiddleware middlewares.MiddlewareFunc
}

func init() {
	Register(new(TodoHandler))
}

func (t *TodoHandler) Init(a *app.App) error {
	t.Router = a.Router
	t.todoService = a.Services.Todo
	t.authMiddleware = middlewares.AuthMiddleware(a.Config.Env, a.Services.Session)
	t.TemplateRenderer = a.TemplateRenderer
	return nil
}

func (t *TodoHandler) Routes() {
	// INFO: to apply middleware to a single route use With method
	// t.Router.With(t.authMiddleware).Get("/todos", t.List)

	t.Route("/todos", func(r infrastructure.Router) {
		// INFO: to apply middleware to a group of routes use Use method
		// r.Use(t.authMiddleware)

		r.Get("/", t.List)
		r.Get("/{id}", t.Get)
		r.Post("/", t.Create)
		r.Put("/{id}", t.Update)
		r.Delete("/{id}", t.Delete)
	})
}

func (t *TodoHandler) List(w http.ResponseWriter, r *http.Request) {
	todos, err := t.todoService.ListTodos(r.Context())

	// NOTE: example of how to return JSON response
	// if the endpoint is handling both HTML and JSON responses
	if t.Wants(r, "application/json") {
		if err != nil {
			jsonErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		jsonResponse(w, http.StatusOK, todos)
		return
	}

	p := page.NewTodoListPage()

	if err != nil {
		p.Error = err.Error()
		w.WriteHeader(500)
		t.Render(w, p)
		return
	}

	todoListPageData := page.NewTodoListPageData(r)
	p.Data = todoListPageData
	todoListPageData.Todos = todos

	w.WriteHeader(200)
	t.Render(w, p)
}

func (t *TodoHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := parseToInt(infrastructure.GetParam(r, "id"))
	if err != nil {
		w.WriteHeader(400)
		w.Header().Add("Hx-Location", "/todos")
		return
	}

	p := page.NewTodoPage(r)

	todo, err := t.todoService.GetTodo(r.Context(), id)
	if err != nil {
		p.Error = err.Error()
		w.WriteHeader(404)
		t.Render(w, p)
		return
	}

	p.Data = page.NewTodoPageData(r, todo)
	w.WriteHeader(200)
	t.Render(w, p)
}

func (t *TodoHandler) Create(w http.ResponseWriter, r *http.Request) {
	input := service.CreateTodoInput{}

	form := page.NewTodoCreateForm(r)
	title := strings.TrimSpace(r.FormValue("title"))
	description := strings.TrimSpace(r.FormValue("description"))
	form.Title = title
	form.Description = description

	if title == "" {
		form.Error = "Title is required"
		w.WriteHeader(400)
		t.RenderPartial(w, "todo-create-form", form)
		return
	}

	input.Title = title
	input.Description = description

	todo, err := t.todoService.CreateTodo(r.Context(), input)
	if err != nil {
		form.Error = err.Error()
		w.WriteHeader(400)
		t.RenderPartial(w, "todo-create-form", form)
		return
	}

	w.WriteHeader(201)
	t.RenderPartial(w, "todo-item-oob", todo)
	t.RenderPartial(w, "todo-create-form", page.NewTodoCreateForm(r))
}

func (t *TodoHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := parseToInt(infrastructure.GetParam(r, "id"))

	if err != nil {
		w.WriteHeader(400)
		w.Header().Add("Hx-Location", "/todos")
		return
	}

	todo, err := t.todoService.GetTodo(ctx, id)
	if err != nil {
		w.WriteHeader(404)
		w.Header().Add("Hx-Location", "/todos")
		return
	}

	todoForm := page.NewTodoEditForm(r, todo)
	title := strings.TrimSpace(r.FormValue("title"))
	description := strings.TrimSpace(r.FormValue("description"))
	todoForm.Title = title
	todoForm.Description = description

	if title == "" {
		todoForm.Error = "Title is required"
		w.WriteHeader(400)
		t.RenderPartial(w, "todo-edit-form", todoForm)
		return
	}

	input := service.UpdateTodoInput{
		Title:       title,
		Description: description,
	}

	todo, err = t.todoService.UpdateTodo(ctx, id, input)
	if err != nil {
		todoForm.Error = err.Error()
		w.WriteHeader(400)
		t.RenderPartial(w, "todo-edit-form", todoForm)
		return
	}

	w.Header().Add("Hx-Trigger", "close_edit_form")
	w.WriteHeader(200)
	t.RenderPartial(w, "todo-details-info-oob", todo)
	t.RenderPartial(w, "todo-edit-form", page.NewTodoEditForm(r, todo))
}

func (t *TodoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := parseToInt(infrastructure.GetParam(r, "id"))

	if err != nil {
		w.WriteHeader(400)
		w.Header().Add("Hx-Location", "/todos")
		return
	}

	todo, err := t.todoService.GetTodo(ctx, id)
	if err != nil {
		w.WriteHeader(404)
		w.Header().Add("Hx-Location", "/todos")
		return
	}

	p := page.NewTodoPage(r)
	templateKey := fmt.Sprintf("%s:%s", p.Layout, p.Name)

	tmpl, err := t.GetTemplate(templateKey)
	if err != nil {
		w.WriteHeader(500)
		t.RenderString(w, err.Error(), nil)
		return
	}

	deleteForm := page.NewDeleteForm(r, todo)

	if err = t.todoService.DeleteTodo(ctx, todo); err != nil {
		w.WriteHeader(500)
		tmplString := fmt.Sprintf("<div id=\"error-message\" hx-swap-oob=\"true\"><p style='color: red;'>%s</p></div>", err.Error())
		tmpl.ExecuteTemplate(w, "todo-delete-form", deleteForm)
		t.RenderString(w, tmplString, nil)
		return
	}

	w.Header().Add("Hx-Location", "/todos")
	w.WriteHeader(200)
}
