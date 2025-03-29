package query

import (
	"go-starter-template/internal/application/result"
	"go-starter-template/internal/domain/entities"
)

type GetTodoQuery struct {
	Todo *result.TodoResult
}

type GetTodoListQuery struct {
	Todos []*result.TodoResult
}

func NewGetTodoQuery(todo *entities.Todo) *GetTodoQuery {
	return &GetTodoQuery{
		Todo: result.NewTodoResult(todo),
	}
}

func NewGetTodoListQuery(todos []*entities.Todo) *GetTodoListQuery {
	var todoResults []*result.TodoResult

	for _, todo := range todos {
		todoResults = append(todoResults, result.NewTodoResult(todo))
	}

	return &GetTodoListQuery{
		Todos: todoResults,
	}
}
