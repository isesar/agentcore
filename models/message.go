package models

import "time"

type Message struct {
	ID             int       `json:"id" db:"id"`
	ConversationID int       `json:"conversation_id" db:"conversation_id"`
	Role           string    `json:"role" db:"role"`
	Content        string    `json:"content" db:"content"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}
