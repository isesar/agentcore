package repository

import "agentcore/models"

type ConversationRepository interface {
	Create(conversation *models.Conversation) (*models.Conversation, error)
	GetByID(id int) (*models.Conversation, error)
	ListByUserID(userID int) ([]models.Conversation, error)
}

type MessageRepository interface {
	Create(message *models.Message) (*models.Message, error)
	GetByID(id int) (*models.Message, error)
	ListByConversationID(conversationID int) ([]models.Message, error)
}

type AgentRunRepository interface {
	Create(run *models.AgentRun) (*models.AgentRun, error)
	GetByID(id int) (*models.AgentRun, error)
	ListByConversationID(conversationID int) ([]models.AgentRun, error)
}

type ToolCallRepository interface {
	Create(toolCall *models.ToolCall) (*models.ToolCall, error)
	ListByAgentRunID(agentRunID int) ([]models.ToolCall, error)
}

type UsageEventRepository interface {
	Create(event *models.UsageEvent) (*models.UsageEvent, error)
	ListByAgentRunID(agentRunID int) ([]models.UsageEvent, error)
}

type ErrorEventRepository interface {
	Create(event *models.ErrorEvent) (*models.ErrorEvent, error)
	ListByAgentRunID(agentRunID int) ([]models.ErrorEvent, error)
}

type RoutingDecisionRepository interface {
	Create(decision *models.RoutingDecisionRecord) (*models.RoutingDecisionRecord, error)
	ListByConversationID(conversationID string) ([]models.RoutingDecisionRecord, error)
	ListByIntent(intent string) ([]models.RoutingDecisionRecord, error)
	ListByTraceID(traceID string) ([]models.RoutingDecisionRecord, error)
}
