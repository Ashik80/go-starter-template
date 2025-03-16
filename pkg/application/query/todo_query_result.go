package query

import (
	"go-starter-template/pkg/application/result"
	"go-starter-template/pkg/domain/entities"
)

type TodoQueryResult struct {
	Todo *result.TodoResult
}

type TodoListQueryResult struct {
	Todos []*result.TodoResult
}

func NewTodoQueryResult(todo *entities.Todo) *TodoQueryResult {
	return &TodoQueryResult{
		Todo: result.NewTodoResult(todo),
	}
}

func NewTodoListQueryResult(todos []*entities.Todo) *TodoListQueryResult {
	var todoResults []*result.TodoResult

	for _, todo := range todos {
		todoResults = append(todoResults, result.NewTodoResult(todo))
	}

	return &TodoListQueryResult{
		Todos: todoResults,
	}
}
