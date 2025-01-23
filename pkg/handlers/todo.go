package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"go-starter-template/pkg/app"
	"go-starter-template/pkg/service"
	"go-starter-template/pkg/store"
)

type TodoHandler struct {
	router    service.Router
	todoStore *store.TodoStore
}

func init() {
	Register(new(TodoHandler))
}

func (t *TodoHandler) Init(a *app.App) error {
	t.router = a.Router
	t.todoStore = a.Store.TodoStore
	return nil
}

func (t *TodoHandler) Routes(router service.Router) {
	router.HandleFunc("/todos", t.List)
	router.HandleFunc("/todos/{id}", t.Get)
	router.HandleFunc("POST /todos", t.Create)
}

func (t *TodoHandler) List(w http.ResponseWriter, r *http.Request) {
	todos, err := t.todoStore.List(r.Context())
	if err != nil {
		jsonErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	jsonResponse(w, http.StatusOK, todos)
}

func (t *TodoHandler) Get(w http.ResponseWriter, r *http.Request) {
	params := t.router.WithPathParams(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		jsonErrorResponse(w, http.StatusBadRequest, fmt.Errorf("invalid id: %v\n", err).Error())
		return
	}
	todo, err := t.todoStore.Get(r.Context(), id)
	if err != nil {
		jsonErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}
	jsonResponse(w, http.StatusOK, todo)
}

func (t *TodoHandler) Create(w http.ResponseWriter, r *http.Request) {
	var todoDto store.TodoCreateDto
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
}
