package models

type QueryContextMetadata struct {
	ConversationID string   `json:"conversation_id,omitempty" binding:"omitempty,max=128"`
	UserID         string   `json:"user_id,omitempty" binding:"omitempty,max=128"`
	TraceID        string   `json:"trace_id,omitempty" binding:"omitempty,max=128"`
	Source         string   `json:"source,omitempty" binding:"omitempty,max=128"`
	Tags           []string `json:"tags,omitempty" binding:"omitempty,dive,max=64"`
}

type QueryRequest struct {
	Prompt  string               `json:"prompt" binding:"required,min=1,max=4000"`
	Context QueryContextMetadata `json:"context,omitempty"`
}

type QueryResponse struct {
	ConversationID string  `json:"conversation_id,omitempty"`
	Answer         string  `json:"answer"`
	Intent         string  `json:"intent"`
	Confidence     float64 `json:"confidence"`
	Reason         string  `json:"reason"`
}

type ErrorPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

type ErrorEnvelope struct {
	Error     ErrorPayload `json:"error"`
	RequestID string       `json:"request_id,omitempty"`
}
