package repository

import "database/sql"

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Users:    NewPQUserStore(db),
		Sessions: NewPQSessionStore(db),
		Todos:    NewPQTodoStore(db),
	}
}
