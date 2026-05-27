package repository

import "database/sql"

type PostgresConversationRepository struct {
	db *sql.DB
}

type PostgresMessageRepository struct {
	db *sql.DB
}

type PostgresAgentRunRepository struct {
	db *sql.DB
}

type PostgresToolCallRepository struct {
	db *sql.DB
}

type PostgresUsageEventRepository struct {
	db *sql.DB
}

type PostgresErrorEventRepository struct {
	db *sql.DB
}

type PostgresRoutingDecisionRepository struct {
	db *sql.DB
}

func NewConversationRepository(db *sql.DB) ConversationRepository {
	return &PostgresConversationRepository{db: db}
}

func NewMessageRepository(db *sql.DB) MessageRepository {
	return &PostgresMessageRepository{db: db}
}

func NewAgentRunRepository(db *sql.DB) AgentRunRepository {
	return &PostgresAgentRunRepository{db: db}
}

func NewToolCallRepository(db *sql.DB) ToolCallRepository {
	return &PostgresToolCallRepository{db: db}
}

func NewUsageEventRepository(db *sql.DB) UsageEventRepository {
	return &PostgresUsageEventRepository{db: db}
}

func NewErrorEventRepository(db *sql.DB) ErrorEventRepository {
	return &PostgresErrorEventRepository{db: db}
}

func NewRoutingDecisionRepository(db *sql.DB) RoutingDecisionRepository {
	return &PostgresRoutingDecisionRepository{db: db}
}
