package store

import "database/sql"

type Store struct {
	db           *sql.DB
	TodoStore    TodoStore
	UserStore    UserStore
	SessionStore SessionStore
}

func NewDataStore(db *sql.DB) *Store {
	s := new(Store)
	s.db = db

	s.initUserStore()
	s.initSessionStore()
	s.initTodoStore()

	return s
}

func (s *Store) initTodoStore() {
	s.TodoStore = NewPQTodoStore(s.db)
}

func (s *Store) initUserStore() {
	s.UserStore = NewPQUserStore(s.db)
}

func (s *Store) initSessionStore() {
	s.SessionStore = NewPQSessionStore(s.db)
}
