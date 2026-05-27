package models

import "time"

type ErrorEvent struct {
	ID         int       `json:"id" db:"id"`
	AgentRunID *int      `json:"agent_run_id,omitempty" db:"agent_run_id"`
	ErrorType  string    `json:"error_type" db:"error_type"`
	Message    string    `json:"message" db:"message"`
	StackTrace string    `json:"stack_trace" db:"stack_trace"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}
