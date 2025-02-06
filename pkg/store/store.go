package store

import (
	"go-starter-template/ent"
)

type Store struct {
	orm          *ent.Client
	TodoStore    TodoStore
	UserStore    UserStore
	SessionStore SessionStore
}

func NewDataStore(orm *ent.Client) *Store {
	s := new(Store)
	s.orm = orm

	s.initUserStore()
	s.initSessionStore()
	s.initTodoStore()

	return s
}

func (s *Store) initTodoStore() {
	s.TodoStore = NewEntTodoStore(s.orm)
}

func (s *Store) initUserStore() {
	s.UserStore = NewEntUserStore(s.orm)
}

func (s *Store) initSessionStore() {
	s.SessionStore = NewEntSessionStore(s.orm)
}
