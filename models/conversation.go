package models

import "time"

type Conversation struct {
	ID        int       `json:"id" db:"id"`
	UserID    *int      `json:"user_id,omitempty" db:"user_id"`
	Title     string    `json:"title" db:"title"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
