package entity

import "time"

type User struct {
	ID        int        `json:"id,omitempty"`
	Email     string     `json:"email,omitempty"`
	Password  string     `json:"password,omitempty"`
	CreatedAt time.Time  `json:"created_at,omitempty"`
	UpdatedAt time.Time  `json:"updated_at,omitempty"`
	Sessions  []*Session `json:"sessions,omitempty"`
}
