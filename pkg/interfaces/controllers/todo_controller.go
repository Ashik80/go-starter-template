package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"go-starter-template/pkg/application/command"
	"go-starter-template/pkg/application/interfaces"
	"go-starter-template/pkg/application/result"
	"go-starter-template/pkg/infrastructure/config"
	"go-starter-template/pkg/infrastructure/csrf"
	"go-starter-template/pkg/infrastructure/middlewares"
	"go-starter-template/pkg/infrastructure/renderer"
	"go-starter-template/pkg/infrastructure/router"
	"go-starter-template/pkg/infrastructure/views/components"
	"go-starter-template/pkg/infrastructure/views/pages"
	"go-starter-template/pkg/utils"
)

type TodoController struct {
	todoService interfaces.TodoService
}

func NewTodoController(r router.Router, todoService interfaces.TodoService, sessionService interfaces.SessionService, config *config.Config) {
	controller := &TodoController{
		todoService: todoService,
	}

	_ = middlewares.AuthMiddleware(config.Env, sessionService)

	r.Route("/todos", func(r router.Router) {
		// INFO: to apply middleware to a group of routes use Use method
		// r.Use(authMiddleware)

		r.Get("/", controller.List)
		r.Get("/{id}", controller.Get)
		r.Post("/", controller.Create)
		r.Put("/{id}", controller.Update)
		r.Delete("/{id}", controller.Delete)
	})
}

func (tc *TodoController) List(w http.ResponseWriter, r *http.Request) {
	data := pages.TodosPageData{}
	data.Form = components.NewTodoCreateForm(r)

	res, err := tc.todoService.ListTodos(r.Context())

	if err != nil {
		data.Error = err.Error()
		w.WriteHeader(500)
		pages.Todos(data).Render(r.Context(), w)
		return
	}

	data.Todos = res.Todos

	w.WriteHeader(200)
	pages.Todos(data).Render(r.Context(), w)
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

	data := pages.TodoDetailsPageData{}

	data.Todo = res.Todo
	data.EditForm = components.NewTodoEditForm(r, res.Todo)
	data.DeleteForm = components.NewTodoDeleteForm(r, res.Todo.ID)

	w.WriteHeader(200)
	pages.TodoDetails(data).Render(r.Context(), w)
}

func (tc *TodoController) Create(w http.ResponseWriter, r *http.Request) {
	form := &components.TodoCreateFormData{
		CSRF:        csrf.GetCSRFField(r),
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
		components.TodoCreateForm(form).Render(r.Context(), w)
		return
	}

	form = components.NewTodoCreateForm(r)

	w.WriteHeader(201)
	components.TodoItemOOB(result.Todo).Render(r.Context(), w)
	components.TodoCreateForm(form).Render(r.Context(), w)
}

func (tc *TodoController) Update(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseToInt(router.GetParam(r, "id"))
	if err != nil {
		w.Header().Add("Hx-Location", "/todos")
		w.WriteHeader(400)
		return
	}

	editForm := components.NewTodoEditForm(r, &result.TodoResult{
		ID:          id,
		Title:       strings.TrimSpace(r.FormValue("title")),
		Description: strings.TrimSpace(r.FormValue("description")),
	})

	result, err := tc.todoService.UpdateTodo(r.Context(), &command.UpdateTodoCommand{
		ID:          id,
		Title:       editForm.Title,
		Description: editForm.Description,
	})

	if err != nil {
		editForm.Error = err.Error()
		w.WriteHeader(400)
		components.TodoEditForm(editForm).Render(r.Context(), w)
		return
	}

	editForm = components.NewTodoEditForm(r, result.Todo)

	w.Header().Add("Hx-Trigger", "close_edit_form")
	w.WriteHeader(200)

	components.TodoDetailsInfoOOB(result.Todo).Render(r.Context(), w)
	components.TodoEditForm(editForm).Render(r.Context(), w)
}

func (tc *TodoController) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseToInt(router.GetParam(r, "id"))
	if err != nil {
		w.Header().Add("Hx-Location", "/todos")
		w.WriteHeader(400)
		return
	}

	deleteForm := components.NewTodoDeleteForm(r, id)
	if err := tc.todoService.DeleteTodo(r.Context(), &command.DeleteTodoCommand{ID: id}); err != nil {
		w.WriteHeader(400)
		components.TodoDeleteForm(deleteForm).Render(r.Context(), w)
		errorString := fmt.Sprintf("<div id=\"error-message\" hx-swap-oob=\"true\"><p style='color: red;'>%s</p></div>", err.Error())
		renderer.RenderString(w, errorString, nil)
		return
	}

	w.Header().Add("Hx-Location", "/todos")
	w.WriteHeader(200)
}
