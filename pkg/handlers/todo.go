package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"go-starter-template/pkg/app"
	"go-starter-template/pkg/entity"
	"go-starter-template/pkg/form"
	"go-starter-template/pkg/infrastructure"
	"go-starter-template/pkg/infrastructure/renderer"
	"go-starter-template/pkg/middlewares"
	"go-starter-template/pkg/service"
	"go-starter-template/pkg/tmpl"
	partialTmpl "go-starter-template/pkg/tmpl/partials"
)

type TodoHandler struct {
	infrastructure.Router
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

	p := tmpl.NewTodosPage()
	pageData := &tmpl.TodosPageData{
		Form:  form.NewTodoCreateForm(r),
		Todos: []*entity.Todo{},
	}

	if err != nil {
		pageData.Error = err.Error()
		w.WriteHeader(500)
		p.Execute(w, pageData)
		return
	}

	pageData.Todos = todos

	w.WriteHeader(200)
	p.Execute(w, pageData)
}

func (t *TodoHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := parseToInt(infrastructure.GetParam(r, "id"))
	if err != nil {
		w.Header().Add("Hx-Location", "/todos")
		w.WriteHeader(400)
		return
	}

	todo, err := t.todoService.GetTodo(r.Context(), id)
	if err != nil {
		w.Header().Add("Hx-Location", "/todos")
		w.WriteHeader(404)
		return
	}

	p := tmpl.NewTodoDetailsPage()
	pageData := &tmpl.TodoDetailsPageData{
		Todo:       todo,
		EditForm:   form.NewTodoEditForm(r, todo),
		DeleteForm: form.NewTodoDeleteForm(r, todo.ID),
	}

	w.WriteHeader(200)
	p.Execute(w, pageData)
}

func (t *TodoHandler) Create(w http.ResponseWriter, r *http.Request) {
	input := service.CreateTodoInput{}

	createForm := form.NewTodoCreateForm(r)
	createForm.Title = strings.TrimSpace(r.FormValue("title"))
	createForm.Description = strings.TrimSpace(r.FormValue("description"))

	createFormTmpl := partialTmpl.NewTodoCreateForm()

	if createForm.Title == "" {
		createForm.Error = "Title is required"
		w.WriteHeader(400)
		createFormTmpl.Execute(w, createForm)
		return
	}

	input.Title = createForm.Title
	input.Description = createForm.Description

	todo, err := t.todoService.CreateTodo(r.Context(), input)
	if err != nil {
		createForm.Error = err.Error()
		w.WriteHeader(400)
		createFormTmpl.Execute(w, createForm)
		return
	}

	todoPartialOobTmpl := partialTmpl.NewTodoItemOob()
	createForm = form.NewTodoCreateForm(r)

	w.WriteHeader(201)
	todoPartialOobTmpl.Execute(w, todo)
	createFormTmpl.Execute(w, createForm)
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

	todoEditForm := form.NewTodoEditForm(r, todo)
	todoEditForm.Title = strings.TrimSpace(r.FormValue("title"))
	todoEditForm.Description = strings.TrimSpace(r.FormValue("description"))

	todoEditFormTmpl := partialTmpl.NewTodoEditForm()

	if todoEditForm.Title == "" {
		todoEditForm.Error = "Title is required"
		w.WriteHeader(400)
		todoEditFormTmpl.Execute(w, todoEditForm)
		return
	}

	input := service.UpdateTodoInput{
		Title:       todoEditForm.Title,
		Description: todoEditForm.Description,
	}

	todo, err = t.todoService.UpdateTodo(ctx, id, input)
	if err != nil {
		todoEditForm.Error = err.Error()
		w.WriteHeader(400)
		todoEditFormTmpl.Execute(w, todoEditForm)
		return
	}

	todoOobTmpl := partialTmpl.NewTodoDetailsInfoOob()
	todoEditForm = form.NewTodoEditForm(r, todo)

	w.Header().Add("Hx-Trigger", "close_edit_form")
	w.WriteHeader(200)

	todoOobTmpl.Execute(w, todo)
	todoEditFormTmpl.Execute(w, todoEditForm)
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

	p := tmpl.NewTodoDetailsPage()
	tmpl := p.GetPageTemplate()

	if err = t.todoService.DeleteTodo(ctx, todo); err != nil {
		deleteForm := form.NewTodoDeleteForm(r, todo.ID)
		tmpl.ExecuteTemplate(w, "todo-delete-form", deleteForm)
		tmplString := fmt.Sprintf("<div id=\"error-message\" hx-swap-oob=\"true\"><p style='color: red;'>%s</p></div>", err.Error())
		renderer.RenderString(w, tmplString, nil)
		return
	}

	w.Header().Add("Hx-Location", "/todos")
	w.WriteHeader(200)
}
