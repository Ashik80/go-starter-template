package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"go-starter-template/pkg/app"
	"go-starter-template/pkg/helpers/todo_helpers"
	"go-starter-template/pkg/middlewares"
	"go-starter-template/pkg/page"
	"go-starter-template/pkg/service"
	"go-starter-template/pkg/store"
)

type TodoHandler struct {
	service.Router
	service.TemplateRenderer
	todoStore      store.TodoStore
	authMiddleware middlewares.MiddlewareFunc
}

func init() {
	Register(new(TodoHandler))
}

func (t *TodoHandler) Init(a *app.App) error {
	t.Router = a.Router
	t.todoStore = a.Store.TodoStore
	t.authMiddleware = middlewares.AuthMiddleware(a.Store.SessionStore)
	t.TemplateRenderer = a.TemplateRenderer
	return nil
}

func (t *TodoHandler) Routes() {
	// INFO: to apply middleware to a single route use With method
	// t.Router.With(t.authMiddleware).Get("/todos", t.List)

	t.Route("/todos", func(r service.Router) {
		// INFO: to apply middleware to a group of routes use Use method
		r.Use(t.authMiddleware)

		r.Get("/", t.List)
		r.Get("/{id}", t.Get)
		r.Post("/", t.Create)
		r.Put("/{id}", t.Update)
		r.Delete("/{id}", t.Delete)
	})
}

func (t *TodoHandler) List(w http.ResponseWriter, r *http.Request) {
	todos, err := t.todoStore.List(r.Context())

	if t.Wants(r, "application/json") {
		if err != nil {
			jsonErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		jsonResponse(w, http.StatusOK, todos)
		return
	}

	todoData := page.NewTodoListPageData(r)
	p := page.NewTodoListPage()

	p.Data = todoData

	if err != nil {
		p.Error = err.Error()
		t.Render(w, p)
		return
	}

	todoData.Todos = todos

	w.WriteHeader(200)
	t.Render(w, p)
}

func (t *TodoHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := parseToInt(service.GetParam(r, "id"))

	if t.Wants(r, "application/json") {
		if err != nil {
			jsonErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		todo, err := t.todoStore.Get(r.Context(), id)
		if err != nil {
			jsonErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		jsonResponse(w, http.StatusOK, todo)
		return
	}

	p := page.NewTodoPage(r)

	todo, err := t.todoStore.Get(r.Context(), id)
	if err != nil {
		p.Error = err.Error()
		w.WriteHeader(404)
		t.Render(w, p)
	}

	todoData := page.NewTodoPageData(r, todo)
	p.Data = todoData

	w.WriteHeader(200)
	t.Render(w, p)
}

func (t *TodoHandler) Create(w http.ResponseWriter, r *http.Request) {
	var todoDto store.TodoCreateDto

	if t.Wants(r, "application/json") {
		err := parseJson(r, &todoDto)
		if err != nil {
			jsonErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		todo, err := t.todoStore.Create(r.Context(), &todoDto)
		if err != nil {
			jsonErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		jsonResponse(w, http.StatusOK, todo)
		return
	}

	title := r.FormValue("title")
	description := r.FormValue("description")

	form := page.NewTodoForm(r)
	form.Title = title
	form.Description = description

	if ok, err := todo_helpers.ValidateTodoForm(form); !ok {
		form.Error = strings.Join(err, ", ")
		w.WriteHeader(400)
		t.RenderPartial(w, "todo-form", form)
		return
	}

	todoDto.Title = title
	todoDto.Description = description

	todo, err := t.todoStore.Create(r.Context(), &todoDto)
	if err != nil {
		form.Error = err.Error()
		w.WriteHeader(400)
		t.RenderPartial(w, "todo-form", form)
		return
	}

	w.WriteHeader(200)
	t.RenderPartial(w, "todo-form", page.NewTodoForm(r))
	t.RenderPartial(w, "todo-item-oob", todo)
}

func (t *TodoHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := parseToInt(service.GetParam(r, "id"))

	var todoDto store.TodoCreateDto

	if t.Wants(r, "application/json") {
		if err != nil {
			jsonErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if err := parseJson(r, &todoDto); err != nil {
			jsonErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		todo, err := t.todoStore.Get(ctx, id)
		if err != nil {
			jsonErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		todo, err = t.todoStore.Update(ctx, todo, &todoDto)
		if err != nil {
			jsonErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		jsonResponse(w, http.StatusOK, todo)
		return
	}

	todo, err := t.todoStore.Get(ctx, id)
	if err != nil {
		todoForm := page.NewTodoEditForm(r, nil)
		todoForm.Error = err.Error()
		w.WriteHeader(400)
		t.RenderPartial(w, "todo-details-edit-form", todoForm)
		return
	}

	todo.Title = r.FormValue("title")
	todo.Description = r.FormValue("description")

	todoForm := page.NewTodoEditForm(r, todo)
	todoDto.Title = todoForm.Title
	todoDto.Description = todoForm.Description

	todo, err = t.todoStore.Update(ctx, todo, &todoDto)
	if err != nil {
		todoForm.Error = err.Error()
		w.WriteHeader(400)
		t.RenderPartial(w, "todo-details-edit-form", todoForm)
		return
	}

	w.Header().Add("Hx-Trigger", "close_edit_form")
	w.WriteHeader(http.StatusOK)

	t.RenderPartial(w, "todo-details-info-oob", todo)
	t.RenderPartial(w, "todo-details-edit-form", page.NewTodoEditForm(r, todo))
}

func (t *TodoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := parseToInt(service.GetParam(r, "id"))

	if err != nil {
		jsonErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if t.Wants(r, "application/json") {
		todo, err := t.todoStore.Get(ctx, id)
		if err != nil {
			jsonErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		if err = t.todoStore.Delete(ctx, todo); err != nil {
			jsonErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		jsonResponse(w, http.StatusOK, map[string]string{"message": "todo deleted successfully"})
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

	todo, err := t.todoStore.Get(ctx, id)
	if err != nil {
		w.WriteHeader(404)
		w.Header().Add("Hx-Location", "/todos")
		return
	}

	deleteForm := page.NewDeleteForm(r, todo)

	if err = t.todoStore.Delete(ctx, todo); err != nil {
		w.WriteHeader(500)
		tmplString := fmt.Sprintf("<div id=\"error-message\" hx-swap-oob=\"true\"><p style='color: red;'>%s</p></div>", err.Error())
		tmpl.ExecuteTemplate(w, "todo-delete-form", deleteForm)
		t.RenderString(w, tmplString, nil)
		return
	}

	w.Header().Add("Hx-Location", "/todos")
	w.WriteHeader(http.StatusNoContent)
}
