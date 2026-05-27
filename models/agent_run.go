package models

import "time"

type AgentRun struct {
	ID             int        `json:"id" db:"id"`
	ConversationID int        `json:"conversation_id" db:"conversation_id"`
	Model          string     `json:"model" db:"model"`
	Status         string     `json:"status" db:"status"`
	StartedAt      time.Time  `json:"started_at" db:"started_at"`
	FinishedAt     *time.Time `json:"finished_at,omitempty" db:"finished_at"`
}
