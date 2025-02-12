package handlers

import (
	"fmt"
	"net/http"

	"go-starter-template/pkg/app"
	"go-starter-template/pkg/form"
	"go-starter-template/pkg/middlewares"
	"go-starter-template/pkg/page"
	"go-starter-template/pkg/service"
	"go-starter-template/pkg/store"
)

type (
	TodoHandler struct {
		service.Router
		*service.TemplateRenderer
		todoStore      store.TodoStore
		authMiddleware middlewares.MiddlewareFunc
	}

	TodoForm struct {
		form.Form
		Title       string
		Description string
		Error       string
	}

	TodoListData struct {
		Todos []*store.Todo
		Form  *TodoForm
	}

	TodoData struct {
		Todo struct {
			form.Form
			ID          int
			Title       string
			Description string
			CreatedAt   string
			UpdatedAt   string
		}
		Error string
	}
)

func newTodoForm(r *http.Request) *TodoForm {
	return &TodoForm{
		Title:       "",
		Description: "",
		Error:       "",
		Form:        form.NewForm(r),
	}
}

func newTodoListData() *TodoListData {
	return &TodoListData{}
}

func newTodoData(r *http.Request) *TodoData {
	t := new(TodoData)
	t.Todo.Form = form.NewForm(r)
	return t
}

func newTodoListPage(t *TodoListData) *page.Page {
	p := page.New()
	p.Title = "Todos"
	p.Name = "todo"
	p.Path = "/todos"
	p.Data = t
	return p
}

func newTodoPage(t *TodoData) *page.Page {
	p := page.New()
	p.Title = "Todo Details"
	p.Name = "todo-details"
	p.Data = t
	return p
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
	// INFO: to protect a route middleware must be called like this
	// t.Handle("/todos", t.authMiddleware(http.HandlerFunc(t.List)))
	t.HandleFunc("/todos", t.List)
	t.HandleFunc("/todos/{id}", t.Get)
	t.HandleFunc("POST /todos", t.Create)
	t.HandleFunc("PUT /todos/{id}", t.Update)
	t.HandleFunc("DELETE /todos/{id}", t.Delete)
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

	todoData := newTodoListData()
	p := newTodoListPage(todoData)

	todoData.Form = newTodoForm(r)

	if err != nil {
		p.Error = err.Error()
		t.Render(w, p)
		return
	}

	todoData.Todos = todos
	t.Render(w, p)
}

func (t *TodoHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := parseToInt(r.PathValue("id"))

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

	todoData := newTodoData(r)
	p := newTodoPage(todoData)

	todo, err := t.todoStore.Get(r.Context(), id)
	if err != nil {
		p.Error = err.Error()
		p.StatusCode = http.StatusNotFound
		t.Render(w, p)
	}

	createdAt := todo.CreatedAt.Format("January 2, 2006 - 3:04PM")
	updatedAt := todo.UpdatedAt.Format("January 2, 2006 - 3:04PM")

	todoData.Todo.ID = todo.ID
	todoData.Todo.Title = todo.Title
	todoData.Todo.Description = todo.Description
	todoData.Todo.CreatedAt = createdAt
	todoData.Todo.UpdatedAt = updatedAt

	t.Render(w, p)
}

func (t *TodoHandler) Create(w http.ResponseWriter, r *http.Request) {
	todoDto := new(store.TodoCreateDto)

	if t.Wants(r, "application/json") {
		err := parseJson(r, &todoDto)
		if err != nil {
			jsonErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		todo, err := t.todoStore.Create(r.Context(), todoDto)
		if err != nil {
			jsonErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		jsonResponse(w, http.StatusOK, todo)
		return
	}

	title := r.FormValue("title")
	description := r.FormValue("description")
	form := newTodoForm(r)
	form.Title = title
	form.Description = description

	if title == "" {
		form.Error = "Title cannot be empty"
		t.RenderPartial(w, http.StatusBadRequest, "todo-form", form)
		return
	}

	todoDto.Title = title
	todoDto.Description = description
	todo, err := t.todoStore.Create(r.Context(), todoDto)
	if err != nil {
		form.Error = err.Error()
		t.RenderPartial(w, http.StatusBadRequest, "todo-form", form)
		return
	}

	form.Title = ""
	form.Description = ""
	t.RenderPartial(w, http.StatusOK, "todo-form", form)
	t.RenderPartial(w, http.StatusOK, "todo-item-oob", todo)
}

func (t *TodoHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := parseToInt(r.PathValue("id"))

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

		todo, err = t.todoStore.Update(ctx, todo, todoDto)
		if err != nil {
			jsonErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		jsonResponse(w, http.StatusOK, todo)
		return
	}

	todoData := newTodoData(r)
	todoData.Todo.ID = id
	todoData.Todo.Title = r.FormValue("title")
	todoData.Todo.Description = r.FormValue("description")

	todoDto.Title = todoData.Todo.Title
	todoDto.Description = todoData.Todo.Description

	todo, err := t.todoStore.Get(ctx, id)
	if err != nil {
		todoData.Error = err.Error()
		t.RenderPartial(w, http.StatusBadRequest, "todo-details-edit-form", todoData)
		return
	}

	todo, err = t.todoStore.Update(ctx, todo, todoDto)
	if err != nil {
		todoData.Error = err.Error()
		t.RenderPartial(w, http.StatusBadRequest, "todo-details-edit-form", todoData)
		return
	}

	todoData.Todo.Title = todo.Title
	todoData.Todo.Description = todo.Description

	createdAt := todo.CreatedAt.Format("January 2, 2006 - 3:04PM")
	updatedAt := todo.UpdatedAt.Format("January 2, 2006 - 3:04PM")

	todoData.Todo.CreatedAt = createdAt
	todoData.Todo.UpdatedAt = updatedAt

	w.Header().Add("Hx-Trigger", "close_edit_form")
	w.WriteHeader(http.StatusOK)
	t.RenderPartial(w, 0, "todo-details-info-oob", todoData.Todo)
	t.RenderPartial(w, 0, "todo-details-edit-form", todoData)
}

func (t *TodoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := parseToInt(r.PathValue("id"))

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

	todo, err := t.todoStore.Get(ctx, id)
	if err != nil {
		tmplString := fmt.Sprintf("<p style='color: red;'>Not found</p>")
		t.RenderString(w, http.StatusNotFound, tmplString, nil)
		return
	}

	if err = t.todoStore.Delete(ctx, todo); err != nil {
		tmplString := fmt.Sprintf("<p style='color: red;'>%s</p>", err.Error())
		t.RenderString(w, http.StatusBadRequest, tmplString, nil)
		return
	}

	w.Header().Add("Hx-Location", "/todos")
}
