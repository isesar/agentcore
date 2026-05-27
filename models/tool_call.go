package models

import "time"

type ToolCall struct {
	ID         int       `json:"id" db:"id"`
	AgentRunID int       `json:"agent_run_id" db:"agent_run_id"`
	Name       string    `json:"name" db:"name"`
	Arguments  string    `json:"arguments" db:"arguments"`
	Result     string    `json:"result" db:"result"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}
