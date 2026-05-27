package models

import "time"

type UsageEvent struct {
	ID               int       `json:"id" db:"id"`
	AgentRunID       int       `json:"agent_run_id" db:"agent_run_id"`
	PromptTokens     int       `json:"prompt_tokens" db:"prompt_tokens"`
	CompletionTokens int       `json:"completion_tokens" db:"completion_tokens"`
	TotalTokens      int       `json:"total_tokens" db:"total_tokens"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
}
