package models

import "time"

type RoutingDecisionRecord struct {
	ID             int       `json:"id" db:"id"`
	ConversationID string    `json:"conversation_id" db:"conversation_id"`
	AgentRunID     *int      `json:"agent_run_id,omitempty" db:"agent_run_id"`
	Intent         string    `json:"intent" db:"intent"`
	Confidence     float64   `json:"confidence" db:"confidence"`
	Reason         string    `json:"reason" db:"reason"`
	Classifier     string    `json:"classifier" db:"classifier"`
	PolicyVersion  string    `json:"policy_version" db:"policy_version"`
	TraceID        string    `json:"trace_id,omitempty" db:"trace_id"`
	UserID         string    `json:"user_id,omitempty" db:"user_id"`
	Source         string    `json:"source,omitempty" db:"source"`
	Tags           []string  `json:"tags,omitempty" db:"tags"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}
