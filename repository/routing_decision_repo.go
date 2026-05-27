package repository

import (
	"encoding/json"

	"agentcore/models"
)

func (r *PostgresRoutingDecisionRepository) Create(decision *models.RoutingDecisionRecord) (*models.RoutingDecisionRecord, error) {
	var created models.RoutingDecisionRecord

	tags, err := json.Marshal(decision.Tags)
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(
		`INSERT INTO routing_decisions (
			conversation_id, agent_run_id, intent, confidence, reason, classifier,
			policy_version, trace_id, user_id, source, tags
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11::jsonb)
		RETURNING id, conversation_id, agent_run_id, intent, confidence, reason, classifier, policy_version, trace_id, user_id, source, tags, created_at`,
		decision.ConversationID,
		decision.AgentRunID,
		decision.Intent,
		decision.Confidence,
		decision.Reason,
		decision.Classifier,
		decision.PolicyVersion,
		decision.TraceID,
		decision.UserID,
		decision.Source,
		string(tags),
	).Scan(
		&created.ID,
		&created.ConversationID,
		&created.AgentRunID,
		&created.Intent,
		&created.Confidence,
		&created.Reason,
		&created.Classifier,
		&created.PolicyVersion,
		&created.TraceID,
		&created.UserID,
		&created.Source,
		&tags,
		&created.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(tags, &created.Tags); err != nil {
		return nil, err
	}
	return &created, nil
}

func (r *PostgresRoutingDecisionRepository) ListByConversationID(conversationID string) ([]models.RoutingDecisionRecord, error) {
	rows, err := r.db.Query(
		`SELECT id, conversation_id, agent_run_id, intent, confidence, reason, classifier, policy_version, trace_id, user_id, source, tags, created_at
		 FROM routing_decisions
		 WHERE conversation_id = $1
		 ORDER BY created_at DESC`,
		conversationID,
	)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	return scanRoutingDecisionRows(rows)
}

func (r *PostgresRoutingDecisionRepository) ListByIntent(intent string) ([]models.RoutingDecisionRecord, error) {
	rows, err := r.db.Query(
		`SELECT id, conversation_id, agent_run_id, intent, confidence, reason, classifier, policy_version, trace_id, user_id, source, tags, created_at
		 FROM routing_decisions
		 WHERE intent = $1
		 ORDER BY created_at DESC`,
		intent,
	)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	return scanRoutingDecisionRows(rows)
}

func (r *PostgresRoutingDecisionRepository) ListByTraceID(traceID string) ([]models.RoutingDecisionRecord, error) {
	rows, err := r.db.Query(
		`SELECT id, conversation_id, agent_run_id, intent, confidence, reason, classifier, policy_version, trace_id, user_id, source, tags, created_at
		 FROM routing_decisions
		 WHERE trace_id = $1
		 ORDER BY created_at DESC`,
		traceID,
	)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	return scanRoutingDecisionRows(rows)
}

func scanRoutingDecisionRows(rows rowScanner) ([]models.RoutingDecisionRecord, error) {
	decisions := make([]models.RoutingDecisionRecord, 0)
	for rows.Next() {
		var d models.RoutingDecisionRecord
		var tags []byte
		if err := rows.Scan(
			&d.ID,
			&d.ConversationID,
			&d.AgentRunID,
			&d.Intent,
			&d.Confidence,
			&d.Reason,
			&d.Classifier,
			&d.PolicyVersion,
			&d.TraceID,
			&d.UserID,
			&d.Source,
			&tags,
			&d.CreatedAt,
		); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(tags, &d.Tags); err != nil {
			return nil, err
		}
		decisions = append(decisions, d)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return decisions, nil
}

type rowScanner interface {
	Next() bool
	Scan(dest ...interface{}) error
	Err() error
}
