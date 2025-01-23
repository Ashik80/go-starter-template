package store

import (
	"gohtmx/ent"
)

type Store struct {
	orm       *ent.Client
	TodoStore *TodoStore
}

func NewDataStore(orm *ent.Client) *Store {
	s := new(Store)
	s.orm = orm

	s.initTodoStore()

	return s
}

func (s *Store) initTodoStore() {
	s.TodoStore = NewTodoStore(s.orm)
}
